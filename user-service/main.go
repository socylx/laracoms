package main

import (
	"fmt"
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	database "github.com/socylx/laracoms/user-service/db"
	"github.com/socylx/laracoms/user-service/handler"
	"github.com/socylx/laracoms/user-service/model"
	pb "github.com/socylx/laracoms/user-service/proto/user"
	repository "github.com/socylx/laracoms/user-service/repo"
	"github.com/socylx/laracoms/user-service/service"
	micro "go-micro.dev/v4"
	"log"
)

func main() {
	// 创建数据库连接，程序退出时断开连接
	db, err := database.CreateConnection()
	defer func() { _ = db.Close() }()

	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// 和 Laravel 数据库迁移类似
	// 每次启动服务时都会检查，如果数据表不存在则创建，已存在检查是否有修改
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.PasswordReset{})

	// 初始化 Repo 实例用于后续数据库操作
	repo := &repository.UserRepository{Db: db.Debug()}
	resetRepo := &repository.PasswordResetRepository{Db: db}
	// 初始化 token service
	token := &service.TokenService{Repo: repo}

	// 以下是 Micro 创建微服务流程

	srv := micro.NewService(
		micro.Name("user"),
		micro.Version("latest"),
	)

	srv.Init()

	// 获取 Broker 实例
	//pubSub := srv.Server().Options().Broker

	// 注册处理器
	_ = pb.RegisterUserServiceHandler(srv.Server(), &handler.UserService{
		Repo:      repo,
		ResetRepo: resetRepo,
		Token:     token,
	})

	// 启动用户服务
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
