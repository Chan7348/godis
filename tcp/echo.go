package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/Chan7348/godis/lib/logger"
	"github.com/Chan7348/godis/lib/sync/atomic"
	"github.com/Chan7348/godis/lib/sync/wait"
)

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(10 * time.Second)
	c.Conn.Close()
	return nil
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		conn.Close()
		return
	}

	client := &EchoClient{
		Conn: conn,
	}

	h.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				h.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		b := []byte(msg)
		_, _ = conn.Write(b)
		client.Waiting.Done()
	}
}

func (h *EchoHandler) Close() error {
	logger.Info("handler shutting down....")
	h.closing.Set(true)
	h.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}
