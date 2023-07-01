package main

import (
	"fmt"
	"github.com/socylx/laracoms/product-service/db"
	"github.com/socylx/laracoms/product-service/handler"
	"github.com/socylx/laracoms/product-service/model"
	pb "github.com/socylx/laracoms/product-service/proto/product"
	"github.com/socylx/laracoms/product-service/repo"
	micro "go-micro.dev/v4"
	"log"
)

func main() {
	// 创建数据库连接，程序退出时断开连接
	database, err := db.CreateConnection()
	defer func() { _ = database.Close() }()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// 数据库迁移
	database.AutoMigrate(&model.Product{})

	// 初始化 Repo 实例用于后续数据库操作
	productRepo := &repo.ProductRepository{Db: database}

	// 以下是 Micro 创建微服务流程
	srv := micro.NewService(
		micro.Name("product"),
		micro.Version("latest"), // 新增接口版本参数
	)
	srv.Init()

	// 注册处理器
	_ = pb.RegisterProductServiceHandler(srv.Server(), &handler.ProductService{
		ProductRepo: productRepo,
	})

	// 启动商品服务
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
