package common

import (
	"encoding/hex"
	"hash/fnv"
	"net"
	//"github.com/golang/protobuf/proto"
)

func genConnID(addr net.Addr) string {
	a := fnv.New32()
	a.Write([]byte(addr.String()))
	return hex.EncodeToString(a.Sum(nil))
}
