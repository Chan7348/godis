package redis

type Connection interface {
	Write([]byte) (int, error)
	Close() error
	RemoteAddress() string

	SetPassword(string)
	GetPassword() string

	// client should keep its subscribing channels
	Subscribe(channel string)
	UnSubscribe(channel string)
	SubscribeCount() int
	GetChannels() []string

	InMultiState() bool
	SetMultiState(bool)
	GetQueuedCmdLine() [][][]byte
	EnqueueCmd([][]byte)
	ClearQueuedCmds()
	GetWatching() map[string]uint32
	AddTxError(err error)
	GetTxErrors() map[string]uint32

	GetDBIndex() int
	SelectDB(int)

	SetSlave()
	IsSlace() bool

	SetMaster()
	IsMaster() bool

	Name() string
}
