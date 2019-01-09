package common

type ConnStatus int
type ServeState int

const ( //连接状态
	ConnStatus_Unknown ConnStatus = iota
	ConnStatus_Run
	ConnStatus_Close
)

const (
	ServeState_Unknown ServeState = iota
	ServeState_Run     ServeState = 1
	ServeState_Close   ServeState = 2
)

const (
	Max_Buf_Size   = 4096
	Max_Conn_Queue = 1024
)
