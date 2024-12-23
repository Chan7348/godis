package asserts

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/Chan7348/godis/interface/redis"
	"github.com/Chan7348/godis/lib/utils"
	"github.com/Chan7348/godis/redis/protocol"
)

func AssertIntReply(t *testing.T, actual redis.Reply, expected int) {
	intResult, ok := actual.(*protocol.IntReply)
	if !ok {
		t.Errorf("expected int protocol, actually %s, %s", actual.ToBytes(), printStack())
		return
	}
	if intResult.Code != int64(expected) {
		t.Errorf("Expected %d, actually %d, %s", expected, intResult.Code, printStack())
	}
}

func AssertIntReplyGreaterThan(t *testing.T, actual redis.Reply, expected int) {
	intResult, ok := actual.(*protocol.IntReply)
	if !ok {
		t.Errorf("expected int protocol, actually %s, %s", actual.ToBytes(), printStack())
		return
	}
	if intResult.Code < int64(expected) {
		t.Errorf("expected %d, actually %d, %s", expected, intResult.Code, printStack())
	}
}

func AssertBulkReply(t *testing.T, actual redis.Reply, expected string) {
	bulkReply, ok := actual.(*protocol.BulkReply)
	if !ok {
		t.Errorf("expected bulk protocol, actually %s, %s", actual.ToBytes(), printStack())
		return
	}
	if !utils.BytesEquals(bulkReply.Arg, []byte(expected)) {
		t.Errorf("expected %s, actually %s, %s", expected, actual.ToBytes(), printStack())
	}
}

func printStack() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("at %s:%d", file, line)
	}
	return ""
}
