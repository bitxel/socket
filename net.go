package socket

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"syscall"
)

const (
	ENV_FD = "SERVER_FD"
)

var (
	inheritFD = make([]uintptr, 0)
	activeFD  = make([]uintptr, 0)
)

// Listen () returns a net.Listener interface
// Fist, it tries to inherit file descriptors which are save in environment variable by the parent process
// If error occurs, it will create a new Listener
func Listen(nettype, laddr string) (net.Listener, error) {
	addr, err := net.ResolveTCPAddr(nettype, laddr)
	if err != nil {
		return nil, err
	}

	for _, fd := range inheritFD {
		l, err := net.FileListener(os.NewFile(fd, "listener"))
		if err != nil {
			fmt.Printf("listen from fd err:%d", fd)
		}
		if addrEqual(addr, l.Addr()) {
			activeFD = append(activeFD, fd)
			return l, nil
		}
	}
	l, err := net.Listen(nettype, laddr)
	if err == nil {
		v := reflect.ValueOf(l).Elem().FieldByName("fd").Elem()
		fd := uintptr(v.FieldByName("sysfd").Int())
		activeFD = append(activeFD, fd)
	}
	return l, err
}

// Fork and run the new binary in the same path
func Fork() {
	setEnvFD()
	forkAndRun()
}

func getEnvFD() {
	strs := strings.Split(os.Getenv(ENV_FD), ",")
	for _, v := range strs {
		fd, err := strconv.Atoi(v)
		if err == nil {
			inheritFD = append(inheritFD, uintptr(fd))
		}
	}
}

func setEnvFD() {
	str := ""
	for _, fd := range activeFD {
		str += fmt.Sprintf("%d,", fd)
	}
	os.Setenv(ENV_FD, str)
}

func addrEqual(addr1, addr2 net.Addr) bool {
	return addr1.Network() == addr2.Network() && addr1.String() == addr2.String()
}

func forkAndRun() (err error) {
	bin := os.Args[0]
	if _, err = os.Stat(bin); err != nil {
		return
	}
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	maxfd := uintptr(0)
	for _, v := range activeFD {
		if v > maxfd {
			maxfd = v
		}
	}
	files := make([]*os.File, maxfd+1)
	files[syscall.Stdin] = os.Stdin
	files[syscall.Stdout] = os.Stdout
	files[syscall.Stderr] = os.Stderr
	for _, v := range activeFD {
		files[v] = os.NewFile(v, fmt.Sprintf("tcp:%d", v))
	}
	p, err := os.StartProcess(bin, os.Args, &os.ProcAttr{
		Dir:   wd,
		Env:   os.Environ(),
		Files: files,
		Sys:   &syscall.SysProcAttr{},
	})
	if err != nil {
		return err
	}
	fmt.Printf("forked pid:%v\n", p.Pid)
	return
}

func handleSig() {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGUSR2)
	<-sigch

	setEnvFD()
	forkAndRun()
}

func init() {
	getEnvFD()
	//go handleSig()
}
