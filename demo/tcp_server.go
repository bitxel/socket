// This is demo of creating a server which can gracefully restart.
// Live demo here: https://asciinema.org/a/TzcYvmTmHMXQ7BP6ViLVOotu1
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"time"

	"github.com/bitxel/socket"
)

var (
	port = "127.0.0.1:12345"
)

func handleConn(conn net.Conn) {
	defer func() {
		fmt.Printf("pid: %d | Conn closed %s\n", os.Getpid(), conn.RemoteAddr())
		conn.Close()
	}()
	fmt.Printf("pid: %d | Accept conn from %s\n", os.Getpid(), conn.RemoteAddr())
	data := make([]byte, 1000)
	conn.Read(data)
	fmt.Printf("pid: %d | Receive: %v", os.Getpid(), string(data))
}

func main() {
	l, err := socket.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("pid: %d | listen on: %s", os.Getpid(), port)

	go func() {
		for {
			conn, err := l.Accept()
			if err == nil {
				go handleConn(conn)
			}
		}
	}()

	// Handle the singal USR2, to fork and run the new binary
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGUSR2)
	<-sigch
	socket.Fork()

	time.Sleep(time.Second * 10)
	log.Printf("server exit, pid: %d", os.Getpid())
}
