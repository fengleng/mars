package service

import (
	"context"

	pb "github.com/fengleng/dddd/client/go/t2"
)

type T2Service struct {
	pb.UnimplementedT2Server
}

func NewT2Service() *T2Service {
	return &T2Service{}
}


func (s *T2Service) CreateT2(ctx context.Context, req *pb.CreateT2Req) (*pb.CreateT2Rsp, error) {
	return &pb.CreateT2Rsp{}, nil
}

func (s *T2Service) UpdateT2(ctx context.Context, req *pb.UpdateT2Req) (*pb.UpdateT2Rsp, error) {
	return &pb.UpdateT2Rsp{}, nil
}

func (s *T2Service) DeleteT2(ctx context.Context, req *pb.DeleteT2Req) (*pb.DeleteT2Rsp, error) {
	return &pb.DeleteT2Rsp{}, nil
}

func (s *T2Service) GetT2(ctx context.Context, req *pb.GetT2Req) (*pb.GetT2Rsp, error) {
	return &pb.GetT2Rsp{}, nil
}

func (s *T2Service) ListT2(ctx context.Context, req *pb.ListT2Req) (*pb.ListT2Rsp, error) {
	return &pb.ListT2Rsp{}, nil
}
