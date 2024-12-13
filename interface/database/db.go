package database

import (
	"time"

	"github.com/Chan7348/godis/interface/redis"
	"github.com/hdt3213/rdb/core"
)

type CmdLine = [][]byte

type DB interface {
	Exec(client redis.Connection, cmdLine [][]byte) redis.Reply
	AfterClientClose(c redis.Connection)
	Close()
	LoadRDB(dec *core.Decoder) error
}

type KeyEventCallback func(dbIndex int, key string, entity *DataEntity)

type DataEntity struct {
	Data interface{}
}

type DBEngine interface {
	DB
	ExecWithLock(connection redis.Connection, cmdLine [][]byte) redis.Reply
	ExecMulti(connection redis.Connection, watching map[string]uint32, cmdLines []CmdLine) redis.Reply
	GetUndoLogs(dbIndex int, cb func(key string, data *DataEntity, expiration *time.Time) bool)
	RWLocks(dbIndex int, writeKeys []string, readKeys []string)
	RWUnLocks(dbIndex int, writeKeys []string, readKeys []string)
	GetDBSize(dbIndex int) (int, int)
	GetEntity(dbIndex int, key string) (*DataEntity, bool)
	GetExpiration(dbIndex int, key string) *time.Time
	SetKeyInsertedCallback(cb KeyEventCallback)
	SetKeyDeletedCallback(cb KeyEventCallback)
}
