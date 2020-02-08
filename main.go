package main

import (
	"flag"
	"fmt"
	"github.com/gtygo/raft/rpc/server"

	"github.com/gtygo/raft/node"
)
var DefaultAddr string

func init() {
	flag.StringVar(&DefaultAddr, "addr", "127.0.0.1:5001", "rpc ip address+port")
}

func main() {

	flag.Parse()
	fmt.Println(DefaultAddr)
	n:=node.NewNode(DefaultAddr)
	go server.StartRpcServer(n)
	n.Loop()
}
