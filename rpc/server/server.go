package server

import (
	"context"
	node2 "github.com/gtygo/raft/node"
	"github.com/gtygo/raft/rpc/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	pb.UnimplementedRaftServiceServer
}

func (s *Server) RequestVote(ctx context.Context, req *pb.RequestVoteReq) (*pb.RequestVoteResp, error) {
	return &pb.RequestVoteResp{
		Term:        req.Term,
		VoteGranted: true,
	}, nil
}

func (s *Server)AppendEntries(ctx context.Context,req *pb.AppendEntriesReq) (*pb.AppendEntriesResp, error){
	return &pb.AppendEntriesResp{
		Term:                 10,
		Success:              true,
	},nil
}


func  StartRpcServer(n *node2.Node) {
	logrus.Infof("raft rpc server start listening at: %s ...", n.Addr)
	lis, err := net.Listen("tcp", n.Addr)
	if err != nil {
		logrus.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterRaftServiceServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		logrus.Fatal(err)
	}

}