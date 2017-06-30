package socket

import (
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListen(t *testing.T) {
	l, err := Listen("tcp", "127.0.0.1:12345")
	assert.Nil(t, err)

	file, err := l.(*net.TCPListener).File()
	assert.Nil(t, err)
	t.Logf("listen fd:%d", file.Fd())

	pid := os.Getpid()
	p, err := os.FindProcess(pid)
	assert.Nil(t, err)
	err = p.Signal(syscall.SIGUSR2)
	assert.Nil(t, err)
	time.Sleep(5 * time.Second)
}
