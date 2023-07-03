package main

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/socylx/laracoms/common/tracer"
	"github.com/socylx/laracoms/common/wrapper/tracer/opentracing/gin2micro"
	pb "github.com/socylx/laracoms/demo-service/proto/demo"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/web"
	"log"
	"os"
)

type Say struct{}

var (
	cli pb.DemoService
)

func (s *Say) Anything(c *gin.Context) {
	log.Print("Received Say.Anything API request")
	c.JSON(200, map[string]string{
		"text": "你好，学院君",
	})
}

func (s *Say) Hello(c *gin.Context) {
	log.Println("Received Say.Hello API request")

	name := c.Param("name")
	ctx, ok := gin2micro.ContextWithSpan(c)
	if ok == false {
		log.Println("get context err")
	}
	response, err := cli.SayHello(ctx, &pb.DemoRequest{
		Name: name,
	})

	if err != nil {
		c.JSON(500, err)
	}

	c.JSON(200, response)
}

func main() {
	var name = "laracom.api.demo"
	// 初始化追踪器
	gin2micro.SetSamplingFrequency(50)
	t, io, err := tracer.NewTracer(name, os.Getenv("MICRO_TRACE_SERVER"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = io.Close() }()
	opentracing.SetGlobalTracer(t)

	// Create service
	service := web.NewService(
		web.Name(name),
	)
	_ = service.Init()

	// setup Demo Server Client
	cli = pb.NewDemoService("demo", client.DefaultClient)

	// Create RESTful handler (using Gin)
	say := new(Say)
	router := gin.Default()
	r := router.Group("/demo")
	r.Use(gin2micro.TracerWrapper)
	r.GET("/hello", say.Anything)
	r.GET("/hello/:name", say.Hello)

	// Register Handler
	service.Handle("/", router)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
