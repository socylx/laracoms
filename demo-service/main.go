package main

import (
	"context"
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	pb "github.com/socylx/laracoms/demo-service/proto/demo"
	micro "go-micro.dev/v4"
	"log"
)

type DemoServiceHandler struct {
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *pb.DemoRequest, rsp *pb.DemoResponse) error {
	rsp.Text = "你好, " + req.Name
	return nil
}
func main() {
	srv := micro.NewService(
		micro.Name("demo"),
	)
	srv.Init()

	_ = pb.RegisterDemoServiceHandler(srv.Server(), &DemoServiceHandler{})
	if err := srv.Run(); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
