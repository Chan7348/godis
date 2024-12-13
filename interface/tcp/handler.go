package tcp

import (
	"context"
	"net"
)

type HandleFunc func(ctx context.Context, connection net.Conn)

type Handler interface {
	Handle(ctx context.Context, connection net.Conn)
	Close() error
}
