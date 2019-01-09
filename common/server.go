package common

import (
	"fmt"
	"net"

	"github.com/golang/protobuf/proto"
)

var serve *ServerConn

type DealFunc func(ctx T, req []byte) (proto.Message, error)

type ServerConn struct {
	addr     string
	state    ServeState
	sessions map[string]*Session
	funcMap  map[int64]DealFunc
}

func init() {
	serve = &ServerConn{}
}

func NewServer(addr string) *ServerConn {
	return &ServerConn{
		addr:     addr,
		sessions: make(map[string]*Session),
		funcMap:  make(map[int64]DealFunc),
	}
}

func (sc *ServerConn) Start() {
	serve = sc //TODO 待优化
	addr, err := net.ResolveTCPAddr("tcp4", sc.addr)
	if err != nil {
		fmt.Println("ResolveTCPAddr error:", err.Error())
		return
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("net.ListenTCP errror:", err.Error())
		return
	}

	for {
		conn, err := l.AcceptTCP()
		if err != nil {

		}

		l.SyscallConn()
		fmt.Println("have a conn:", conn.RemoteAddr().String())
		connID := genConnID(conn.RemoteAddr())
		if _, ok := sc.sessions[connID]; !ok {
			sc.sessions[connID] = NewSession(sc, conn, connID)
		}
		go sc.sessions[connID].launch()

		// 生成hash值作为Key
		//sc.conns[conn.RemoteAddr().String()] = conn
		//开始监听事件
	}
}

func (sc *ServerConn) RegisterFunc(id int64, fn DealFunc) {
	//fmt.Printf("====>>msg: %+v\n", reflect.TypeOf(msg))
	if _, ok := sc.funcMap[id]; !ok {
		sc.funcMap[id] = fn
	}
}
