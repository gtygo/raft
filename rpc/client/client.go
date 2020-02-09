package client

import (
	"context"
	"github.com/gtygo/raft/rpc/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"time"
)

type Client struct {
	Deadline      time.Duration
	ReqVoteCh     chan uint64
	ReqVoteDoneCh chan struct{}

	AppendEntriesCh     chan uint64
	AppendEntriesDoneCh chan struct{}
}

func (client *Client) DoRequestVote(sendAddr string, req *pb.RequestVoteReq) error {

	conn, err := grpc.Dial(sendAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
		//todo：处理连接失败的情况
	}
	defer conn.Close()

	c := pb.NewRaftServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.RequestVote(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("get resp:",r.Term,r.VoteGranted)
	logrus.Info("start insert into vote channel")
	if r.VoteGranted {
		client.ReqVoteCh <- r.Term
	}
	return nil
}

func (client *Client) DoAppendEntries(sendAddr string, req *pb.AppendEntriesReq) error {
	conn, err := grpc.Dial(sendAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewRaftServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.AppendEntries(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	if r.Success {
		client.AppendEntriesCh <- r.Term
	}

	return nil
}
