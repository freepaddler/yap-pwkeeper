package main

import (
	"context"
	"fmt"

	"github.com/alecthomas/kingpin/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"yap-pwkeeper/internal/pkg/grpc/proto"
)

var (
	user, password, jwt string
	conn                *grpc.ClientConn
	auth                proto.AuthClient
)

func main() {
	var command string
	kingpin.Flag("user", "user").Short('u').StringVar(&user)
	kingpin.Flag("password", "password").Short('p').StringVar(&password)
	kingpin.Flag("jwt", "token").Short('j').StringVar(&jwt)
	kingpin.Arg("command", "command").Required().StringVar(&command)
	kingpin.Parse()
	var err error
	conn, err = grpc.Dial(
		"127.0.0.1:3200",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	auth = proto.NewAuthClient(conn)

	switch command {
	case "login":
		login()
	case "register":
		register()
	case "refresh":
		refresh()
	}
}

func login() {
	fmt.Println("LOGIN")
	req := &proto.LoginCredentials{
		Login:    user,
		Password: password,
	}
	token, err := auth.Login(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.Token)
}

func register() {
	fmt.Println("REGISTER")
	req := &proto.LoginCredentials{
		Login:    user,
		Password: password,
	}
	token, err := auth.Register(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.Token)
}

func refresh() {
	fmt.Println("REFRESH")
	req := &proto.Token{Token: jwt}
	token, err := auth.Refresh(context.Background(), req)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.Token)
}
