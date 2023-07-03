package main

import (
	"context"
	"github.com/gin-gonic/gin"
	pb "github.com/socylx/laracoms/demo-service/proto/demo"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/web"
	"log"
)

type Say struct{}

var (
	cli pb.DemoService
)

func (s *Say) Anything(c *gin.Context) {
	log.Print("Received Say.Anything API request")
	c.JSON(200, map[string]string{
		"message": "你好，学院君",
	})
}

func (s *Say) Hello(c *gin.Context) {
	log.Print("Received Say.Hello API request")

	name := c.Param("name")

	response, err := cli.SayHello(context.TODO(), &pb.DemoRequest{
		Name: name,
	})

	if err != nil {
		c.JSON(500, err)
	}

	c.JSON(200, response)
}

func main() {
	// Create service
	service := web.NewService(
		web.Name("laracom.api.demo"),
	)
	_ = service.Init()

	// setup Demo Server Client
	cli = pb.NewDemoService("demo", client.DefaultClient)

	// Create RESTful handler (using Gin)
	say := new(Say)
	router := gin.Default()
	router.GET("/hello", say.Anything)
	router.GET("/hello/:name", say.Hello)

	// Register Handler
	service.Handle("/", router)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
