package protocol

import (
	"bytes"

	"github.com/Chan7348/godis/interface/redis"
)

type PongReply struct{}

var pongBytes = []byte("+PONG\r\n")

func (*PongReply) ToBytes() []byte {
	return pongBytes
}

type OkReply struct{}

var okBytes = []byte("+OK\r\n")

func (*OkReply) ToBytes() []byte {
	return okBytes
}

var theOkReply = new(OkReply)

func MakeOkReply() *OkReply {
	return theOkReply
}

var nullBulkBytes = []byte("$-1\r\n")

type NullBulkReply struct{}

func (*NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

var emptyMultiBulkBytes = []byte("*0\r\n")

type EmptyMultiBulkReply struct{}

func (*EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

func IsEmptyMultiBulkReply(reply redis.Reply) bool {
	return bytes.Equal(reply.ToBytes(), emptyMultiBulkBytes)
}

type NoReply struct{}

var noBytes = []byte("")

func (*NoReply) ToBytes() []byte {
	return noBytes
}

type QueuedReply struct{}

var queuedBytes = []byte("+QUEUED\r\n")

func (*QueuedReply) ToBytes() []byte {
	return queuedBytes
}

var theQueuedReply = new(QueuedReply)

func MakeQueuedReply() *QueuedReply {
	return theQueuedReply
}
