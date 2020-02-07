package server

import (
	"context"
	"github.com/gtygo/raft/rpc/pb"
	"log"
)

type Server struct {
	pb.UnimplementedRaftServiceServer
}

func (s *Server)RequestVote(ctx context.Context,req *pb.RequestVoteReq)(*pb.RequestVoteResp, error) {
	log.Println("start request vote")
	log.Println("end request vote")
	return &pb.RequestVoteResp{
		Term:                 req.Term,
		VoteGranted:          true,
	},nil
}