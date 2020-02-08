package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gtygo/raft/rpc/client"
	"github.com/gtygo/raft/rpc/pb"
	"github.com/gtygo/raft/rpc/server"
	"google.golang.org/grpc"
)

const (
	FOLLOWER=iota
	CANDIDATE
	LEADER

	DefaultAddr ="localhost:5001"
)
type State int

type Instruction struct {
	Type string
	Key string
	Value string
}

type LogInfo struct {
	Ins Instruction
	Term uint64
}

type Node struct{
	Id uint64
	Addr string
	Time time.Duration
	NodeState State
	server server.Server
	client client.Client


	//Persistent state on all servers:
	CurrentTerm uint64
	VotedFor uint64
	Log []LogInfo

	//Volatile state on all servers:
	CommitIndex uint64
	LastApplied uint64

	//Volatile state on leaders:
	NextIndex []uint64
	MatchIndex []uint64


}


func (n *Node)StartRpcServer(){

	lis,err:=net.Listen("tcp",n.Addr)
	if err!=nil{
		log.Fatal(err)
	}
	s:=grpc.NewServer()
	pb.RegisterRaftServiceServer(s,&n.server)
	fmt.Println("-------got here")
	if err:=s.Serve(lis);err!=nil{
		log.Fatal(err)
	}

}

func Loop(){
	n:=&Node{
		Id:          0,
		Addr:        DefaultAddr,
		NodeState:   FOLLOWER,
		server:      server.Server{},
		client: client.Client{Addr:DefaultAddr,Ch:make(chan uint64,10),DoneCh:make(chan struct{},1)},

	}
	go n.StartRpcServer()
	go n.receiveVote()

	eachTermTime := time.Second * 5
	t:=time.NewTicker(eachTermTime)
	defer t.Stop()

	for {
		<-t.C
		fmt.Println("got new round ,start request vote")
		for i:=0;i<10;i++{
			go n.client.DoRequestVote(&pb.RequestVoteReq{
				Term:                 uint64(i),
				CandidateId:          2,
				LastLogIndex:         3,
				LastLogTerm:          9,
			})

		}
	}
}




func main(){

	Loop()

}

func (n *Node)receiveVote(){
	count:=0
	for {
		select{
		case x:=<-n.client.ReqVoteCh :
			count++
			fmt.Println("received: ",x)
		case <-n.client.ReqVoteDoneCh :
			break
		}
		fmt.Println(count)
	}
}

func (n *Node)receiveAppendEntries(){
	count:=0

	for{
		select {
			case x:=<-n.client.AppendEntriesCh:
				count++
				fmt.Println("received append entries:",x)
			case <-n.client.AppendEntriesDoneCh:
				break
		}
	}
}
