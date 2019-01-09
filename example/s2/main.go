package main

import (
	"coral/apidoc"
	"coral/common"
	"fmt"

	"github.com/golang/protobuf/proto"
)

func main() {
	serve := common.NewServer(":9090")
	serve.RegisterFunc(100, GetName)

	serve.Start()
}

func GetName(ctx common.T, req []byte) (proto.Message, error) {
	fmt.Printf("===>>todo111:%s\n", string(req))

	return &apidoc.Response{
		Reply: "reply",
	}, nil
}
