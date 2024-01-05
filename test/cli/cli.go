package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alecthomas/kingpin/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"yap-pwkeeper/internal/pkg/grpc/proto"
)

var (
	user, password, jwt string
	conn                *grpc.ClientConn
	auth                proto.AuthClient
	wallet              proto.WalletClient
	name, text          string
	meta1               proto.Meta = proto.Meta{
		Key:   "key1",
		Value: "value1",
	}
	meta2 proto.Meta = proto.Meta{
		Key:   "key2",
		Value: "value2",
	}
)

func main() {
	var command string
	kingpin.Flag("user", "user").Short('u').StringVar(&user)
	kingpin.Flag("password", "password").Short('p').StringVar(&password)
	kingpin.Flag("jwt", "token").Short('j').StringVar(&jwt)
	kingpin.Flag("name", "document name").Short('n').StringVar(&name)
	kingpin.Flag("text", "text").Short('t').StringVar(&text)
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
	wallet = proto.NewWalletClient(conn)

	switch command {
	case "login":
		login()
	case "register":
		register()
	case "refresh":
		refresh()
	case "addnote":
		addNote()
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

func addNote() {
	fmt.Println("ADD NOTE")

	req := &proto.Note{
		Name:     name,
		Metadata: []*proto.Meta{&meta1, &meta2},
		Text:     text,
	}
	_, err := wallet.AddNote(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
}
