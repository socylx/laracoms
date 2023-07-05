package main

import (
	"context"
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	traceplugin "github.com/go-micro/plugins/v4/wrapper/trace/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/socylx/laracoms/common/tracer"
	"github.com/socylx/laracoms/common/wrapper/breaker/hystrix"
	pb "github.com/socylx/laracoms/demo-service/proto/demo"
	"go-micro.dev/v4"
	"go-micro.dev/v4/metadata"
	"log"
	"os"
	"time"
)

func main() {
	// 初始化追踪器
	t, io, err := tracer.NewTracer("laracom.demo.cli", os.Getenv("MICRO_TRACE_SERVER"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = io.Close() }()

	hystrix.Configure([]string{"laracom.service.demo.DemoService.SayHello"})
	srv := micro.NewService(
		micro.Name("laracom.demo.cli"),
		micro.WrapClient(traceplugin.NewClientWrapper(t)),
		micro.WrapClient(hystrix.NewClientWrapper()),
	)
	srv.Init()

	client := pb.NewDemoService("demo", srv.Client())

	// 创建空的上下文, 生成追踪 span
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "call")
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}
	defer span.Finish()

	// 注入 opentracing textmap 到空的上下文用于追踪
	opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md))
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = metadata.NewContext(ctx, md)
	// 记录请求 && 响应 && 错误
	req := &pb.DemoRequest{Name: "学院君"}
	span.SetTag("req", req)

	resp, err := client.SayHello(ctx, req)
	if err != nil {
		span.SetTag("err", err)
		log.Fatalf("服务调用失败：%v", err)
		return
	}
	span.SetTag("resp", resp)
	log.Println(resp.Text)

	// let's delay the process for exiting for reasons you'll see below
	time.Sleep(time.Second * 5)
}
