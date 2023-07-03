package main

import (
	"context"
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	traceplugin "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/socylx/laracoms/common/tracer"
	pb "github.com/socylx/laracoms/demo-service/proto/demo"
	micro "go-micro.dev/v4"
	"go-micro.dev/v4/metadata"
	"log"
	"os"
)

type DemoServiceHandler struct {
}

func (s *DemoServiceHandler) SayHello(ctx context.Context, req *pb.DemoRequest, rsp *pb.DemoResponse) error {
	// 从微服务上下文中获取追踪信息
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	var sp opentracing.Span
	wireContext, _ := opentracing.GlobalTracer().Extract(opentracing.TextMap, opentracing.TextMapCarrier(md))
	// 创建新的 Span 并将其绑定到微服务上下文
	sp = opentracing.StartSpan("SayHello", opentracing.ChildOf(wireContext))
	// 记录请求
	sp.SetTag("req", req)
	defer func() {
		// 记录响应
		sp.SetTag("res", rsp)
		// 在函数返回 stop span 之前，统计函数执行时间
		sp.Finish()
	}()

	rsp.Text = "你好, " + req.Name
	return nil
}

func main() {
	// 初始化全局服务追踪
	t, io, err := tracer.NewTracer("demo", os.Getenv("MICRO_TRACE_SERVER"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = io.Close() }()
	opentracing.SetGlobalTracer(t)

	srv := micro.NewService(
		micro.Name("demo"),
		micro.WrapHandler(traceplugin.NewHandlerWrapper(opentracing.GlobalTracer())), // 基于 jaeger 采集追踪数据
	)
	srv.Init()

	_ = pb.RegisterDemoServiceHandler(srv.Server(), &DemoServiceHandler{})
	if err := srv.Run(); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
