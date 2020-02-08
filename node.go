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
	ClusterNodeAddr []string
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
		Id:        0,
		Addr:      DefaultAddr,
		NodeState: FOLLOWER,
		server:    server.Server{},
		client: client.Client{
			ReqVoteCh:           make(chan uint64,10),
			ReqVoteDoneCh:       make(chan struct{},1),
			AppendEntriesCh:     make(chan uint64,10),
			AppendEntriesDoneCh: make(chan struct{},1),
		},
	}
	go n.StartRpcServer()
	go n.receiveVote()

	eachTermTime := time.Second * 5
	t:=time.NewTicker(eachTermTime)
	defer t.Stop()

	for {
		<-t.C
		fmt.Println("got new round ,start request vote")
		// leader election:

		//选举超时时间





		// 选举超时后，变为candidate 进行参选 ，此节点进入参选状态
		n.NodeState=CANDIDATE

		// 1.自增当前任期号
		n.CurrentTerm++
		//todo： 2.为自己投票

		//todo： 3.重置选举超时计时器

		//4. 发送请求投票RPC给所有服务器
		for i:=0;i<len(n.ClusterNodeAddr);i++{
			go n.client.DoRequestVote(n.ClusterNodeAddr[i],&pb.RequestVoteReq{
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
			if  x > n.CurrentTerm{
				n.CurrentTerm=x
				n.NodeState=FOLLOWER
			}



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
