package common

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type FuncResponse struct {
	msg proto.Message
	err error
}

//Session 会话
type Session struct {
	conn     net.Conn
	status   ConnStatus //状态
	lasttime int64      //最近消息
	total    int64      //总数
	start    time.Time  //开始
	bufPool  sync.Pool  //
	recvList chan string
	sendList chan string
	serve    *ServerConn
	reqId    string
}

func NewSession(sc *ServerConn, conn net.Conn, reqId string) *Session {
	return &Session{
		serve: sc,
		conn:  conn,
		reqId: reqId,
		bufPool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, Max_Buf_Size)
				return &b
			},
		},
		recvList: make(chan string, Max_Conn_Queue),
		sendList: make(chan string, Max_Conn_Queue),
	}
}

func (s *Session) launch() {
	go s.process()
	go s.recvLoop()
	go s.sendLoop()
}

func (s *Session) recvLoop() {
	// 接收消息
	// 解析消息头部
	// 查找字典
	// 接口回调
	// 等待返回
	// 接口响应

	//TODO
	//buf := make([]byte, 4096)
	buf := s.bufPool.Get().(*[]byte)
	for {
		runtime.Gosched()
		len, err := s.conn.Read(*buf)
		//fmt.Printf("recv: %d|%+v\n", length, err)
		if err != nil {
			fmt.Printf("read buf failed:%s\n", err.Error())
			return
		}

		s.recvList <- fmt.Sprintf("%s", string((*buf)[0:len]))

		//fmt.Println("recv:", string((*buf)[0:len]))
		//s.bufPool.Put(buf)
	}
}

func (s *Session) process() {
	//映射请求
	//调用接口
	// resp, err := todo(ctx, req)
	// 打包消息
	for {
		runtime.Gosched()

		select {
		case msg := <-s.recvList:
			go s.invoke([]byte(msg))
		}
	}
}

func (s *Session) invoke(data []byte) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ctx = context.WithValue(ctx, "reqId", s.reqId)
	c := make(chan FuncResponse)

	//time.Sleep(time.Second * 1)

	go func(ctx context.Context, c chan FuncResponse) {
		f, ok := s.serve.funcMap[100]
		if !ok {
			c <- FuncResponse{
				msg: nil,
				err: fmt.Errorf("not find function"),
			}
		}

		resp, err := f(ctx, data)
		if err != nil {
			fmt.Println("====>>>resp err:", time.Since(start))
			c <- FuncResponse{
				msg: nil,
				err: err,
			}
			return
		}
		fmt.Println("===>>>resp:", resp, time.Since(start))
		c <- FuncResponse{
			msg: resp,
			err: nil,
		}
	}(ctx, c)

	select {
	case res := <-c:
		if res.err != nil {
			fmt.Println("=====>>>Unexpected success failed!", res)
			return
		}
		fmt.Println("=====>>>Unexpected success!", res.msg.String())
		data, _ := proto.Marshal(res.msg)
		s.conn.Write(data)
	case <-ctx.Done():
		fmt.Println("=====>>invoke done:", ctx.Err(), time.Since(start))
	}
}

func (s *Session) sendLoop() {

}
