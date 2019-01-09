package common

import (
	"context"

	"github.com/golang/protobuf/proto"
)

type T struct {
	context.Context
	msg   proto.Message
	msgId int32
	reqId string
}

func NewContext(ctx context.Context) T {
	return T{
		Context: ctx,
	}
}

func (c T) Write(reqId int64, data []byte) {

}

func (c T) SetMessage(id int32, msg proto.Message) {

}

//Broadcast 广播
func (c T) Broadcast(reqIds []string) {
	for _, reqId := range reqIds {
		if sn, ok := serve.sessions[reqId]; ok {
			sn.sendList <- "broadcast"
		}
	}
}

//Unicast 点播
func (c T) Unicast() {

}

//Global 全局
func (c T) Global() {

}
