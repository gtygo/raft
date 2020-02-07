package client

import (
	"context"
	"github.com/gtygo/raft/rpc/pb"
	"google.golang.org/grpc"
	"log"
	"time"
)


type Client struct{
	Addr string
	Deadline time.Duration
	Ch  chan uint64
	DoneCh chan struct{}
}

func (client *Client)DoRequestVote(req *pb.RequestVoteReq) error{

	conn,err:=grpc.Dial(client.Addr,grpc.WithInsecure(),grpc.WithBlock())
	if err!=nil{
		log.Fatal(err)
	}
	defer conn.Close()

	c:=pb.NewRaftServiceClient(conn)
	ctx,cancel:=context.WithTimeout(context.Background(),time.Second)
	defer cancel()
	r,err:=c.RequestVote(ctx,req)
	if err!=nil{
		log.Fatal(err)
	}
	log.Println("get resp:",r.Term,r.VoteGranted)

	if r.VoteGranted{
		client.Ch <- r.Term
	}

	return nil
}