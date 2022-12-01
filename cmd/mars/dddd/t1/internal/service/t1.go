package service

import (
	"context"
	"github.com/fengleng/mars/log"
	"io"

	pb "github.com/fengleng/dddd/client/go/t1"
)

type T1Service struct {
	pb.UnimplementedT1Server
}

func NewT1Service() *T1Service {
	return &T1Service{}
}


func (s *T1Service) CreateT1(ctx context.Context, req *pb.CreateT1Req) (*pb.CreateT1Rsp, error) {
	log.Info("msg")
	return &pb.CreateT1Rsp{}, nil
}

func (s *T1Service) UpdateT1(ctx context.Context, req *pb.UpdateT1Req) (*pb.UpdateT1Rsp, error) {
	return &pb.UpdateT1Rsp{}, nil
}

func (s *T1Service) DeleteT1(ctx context.Context, req *pb.DeleteT1Req) (*pb.DeleteT1Rsp, error) {
	return &pb.DeleteT1Rsp{}, nil
}

func (s *T1Service) GetT1(ctx context.Context, req *pb.GetT1Req) (*pb.GetT1Rsp, error) {
	return &pb.GetT1Rsp{}, nil
}
func (s *T1Service) ListT1(conn pb.T1_ListT1Server) error {
	for {
		req, err := conn.Recv()
		if err == io.EOF {
			return conn.SendAndClose(&pb.ListT1Rsp{})
		}
		log.Info(req)
		if err != nil {
			return err
		}
	}
}
