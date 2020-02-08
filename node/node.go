package node

import (
	"math/rand"
	"time"

	"github.com/gtygo/raft/rpc/client"
	"github.com/gtygo/raft/rpc/pb"

	"github.com/sirupsen/logrus"
)

const (
	FOLLOWER = iota
	CANDIDATE
	LEADER
)

type State int

type Instruction struct {
	Type  string
	Key   string
	Value string
}

type LogInfo struct {
	Ins  Instruction
	Term uint64
}

type Node struct {
	Id              uint64
	Addr            string
	Time            time.Duration
	NodeState       State
	ClusterNodeAddr []string
	client          client.Client

	//Persistent state on all servers:
	CurrentTerm uint64
	VotedFor    uint64
	Log         []LogInfo

	//Volatile state on all servers:
	CommitIndex uint64
	LastApplied uint64

	//Volatile state on leaders:
	NextIndex  []uint64
	MatchIndex []uint64
}


func NewNode(addr string)*Node{

	n := &Node{
		Id:              0,
		Addr:            addr,
		NodeState:       FOLLOWER,
		ClusterNodeAddr: []string{"127.0.0.1:5001", "127.0.0.1:5002", "127.0.0.1:5003"},
		client: client.Client{
			ReqVoteCh:           make(chan uint64, 10),
			ReqVoteDoneCh:       make(chan struct{}, 1),
			AppendEntriesCh:     make(chan uint64, 10),
			AppendEntriesDoneCh: make(chan struct{}, 1),
		},
	}

	return n
}



func (n *Node)Loop() {


	logrus.Infof("startup raft node: %v ,node state is %v (follower)", n.Id, n.NodeState)

	for {

		switch n.NodeState {

		case FOLLOWER:

			select {
			case x := <-n.client.ReqVoteCh:
				logrus.Infof("follower received request vote msg %v", x)
				n.checkCurrentTerm(x)
			case x := <-n.client.AppendEntriesCh:
				logrus.Infof("follower received append entries msg %v", x)
				n.checkCurrentTerm(x)
			case <-time.After(time.Duration(Random(150, 300)) * time.Millisecond):
				logrus.Infof("follower timeout ,start new round,%v", 1)
				//todo: start leader election
				n.startLeaderElection()
			}

		case CANDIDATE:
			select {
			case x := <-n.client.ReqVoteCh:
				logrus.Infof("candidate received request vote msg %v", x)
				n.checkCurrentTerm(x)
			case x := <-n.client.AppendEntriesCh:
				logrus.Infof("candidate received append entries msg %v", x)
				n.checkCurrentTerm(x)

			case <-time.After(time.Duration(Random(150, 300)) * time.Millisecond):
				logrus.Infof("candidate timeout ,start new round,%v", 1)
				//todo: start leader election

			}

		case LEADER:
			// todo： 接收客户端请求

			n.boardCastHeartBeat()

		}
	}
}

func (n *Node)checkCurrentTerm(T uint64){
	if T>n.CurrentTerm{
		n.CurrentTerm=T
		n.NodeState=FOLLOWER
	}
}

func (n *Node)startLeaderElection(){

}



func (n *Node) receiveVote() {
	logrus.Infof("start listening request vote rpc message....")
	count := 0
	for {
		select {
		case x := <-n.client.ReqVoteCh:
			count++
			handleRequestVoteMessage(x)

		case <-n.client.ReqVoteDoneCh:
			count = 0
			handleRequestVoteDoneMessage()
		}
	}

}

func (n *Node) receiveAppendEntries() {
	logrus.Infof("start listening append entries rpc message....")
	count := 0
	for {
		select {
		case x := <-n.client.AppendEntriesCh:
			handleAppendEntriesMessage(x)
			count++
		case <-n.client.AppendEntriesDoneCh:
			break
		}
	}
}


//不断发送心跳给其他服务器
func (n *Node) boardCastHeartBeat() {
	for i := 0; i < len(n.ClusterNodeAddr); i++ {
		if n.ClusterNodeAddr[i]==n.Addr{
			continue
		}
		time.Sleep(time.Millisecond * 100)
		go n.client.DoAppendEntries(n.ClusterNodeAddr[i], &pb.AppendEntriesReq{})
	}
}

func handleRequestVoteMessage(term uint64) {

}

func handleRequestVoteDoneMessage() {

}

func handleAppendEntriesMessage(term uint64) {

}

func Random(l int, r int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(r-l) + r
}
