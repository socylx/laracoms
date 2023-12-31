package main

import (
	"context"
	_ "github.com/go-micro/plugins/v4/registry/etcd"
	pb "github.com/socylx/laracoms/user-service/proto/user"
	cli "github.com/urfave/cli/v2"
	micro "go-micro.dev/v4"
	"log"
)

func main() {

	// 初始化客户端服务，定义命令行参数标识
	srv := micro.NewService(micro.Flags(
		&cli.StringFlag{
			Name:  "name",
			Usage: "Your Name",
		},
		&cli.StringFlag{
			Name:  "email",
			Usage: "Your Email",
		},
		&cli.StringFlag{
			Name:  "password",
			Usage: "Your Password",
		},
	))
	client := pb.NewUserService("laracom.service.user", srv.Client())

	// 运行客户端命令调用远程服务逻辑设置
	srv.Init(
		micro.Action(func(c *cli.Context) error {
			name := c.String("name")
			email := c.String("email")
			password := c.String("password")

			log.Println("参数:", name, email, password)

			// 调用用户服务
			r, err := client.Create(context.TODO(), &pb.User{
				Name:     name,
				Email:    email,
				Password: password,
			})
			if err != nil {
				log.Fatalf("创建用户失败: %v", err)
			}
			log.Printf("创建用户成功: %s", r.User.Id)

			// 调用用户认证服务
			var token *pb.Token
			token, err = client.Auth(context.TODO(), &pb.User{
				Email:    email,
				Password: password,
			})
			if err != nil {
				log.Fatalf("用户登录失败: %v", err)
			}
			log.Printf("用户登录成功：%s", token.Token)

			// 调用用户验证服务
			token, err = client.ValidateToken(context.TODO(), token)
			if err != nil {
				log.Fatalf("用户认证失败: %v", err)
			}
			log.Printf("用户认证成功：%s", token.Valid)

			getAll, err := client.GetAll(context.Background(), &pb.Request{})
			if err != nil {
				log.Fatalf("获取所有用户失败: %v", err)
			}
			for _, v := range getAll.Users {
				log.Println(v)
			}

			resp, err := client.CreatePasswordReset(context.Background(), &pb.PasswordReset{Email: email, Token: "password_reset_token"})
			if err != nil {
				log.Fatalf("修改密码失败: %v", err)
			}
			log.Println(resp)

			return nil
		}),
	)

	if err := srv.Run(); err != nil {
		log.Fatalf("用户客户端启动失败: %v", err)
	}
}
