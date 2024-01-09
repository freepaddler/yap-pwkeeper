package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"yap-pwkeeper/internal/pkg/grpc/proto"
)

var (
	user, password, jwt            string
	conn                           *grpc.ClientConn
	auth                           proto.AuthClient
	docs                           proto.DocsClient
	name, text, id, expires, fname string
	serial                         int64
	meta1                          proto.Meta = proto.Meta{
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
	kingpin.Flag("jwt", "token").Envar("TOKEN").Short('j').StringVar(&jwt)
	kingpin.Flag("name", "document name").Short('n').StringVar(&name)
	kingpin.Flag("text", "text").Short('t').StringVar(&text)
	kingpin.Flag("serial", "serial").Short('s').Int64Var(&serial)
	kingpin.Flag("id", "id").Short('i').StringVar(&id)
	kingpin.Flag("expires", "expires").Short('e').StringVar(&expires)
	kingpin.Flag("file", "file").Short('f').StringVar(&fname)
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
	defer func() { _ = conn.Close() }()
	auth = proto.NewAuthClient(conn)
	docs = proto.NewDocsClient(conn)

	switch command {
	case "login":
		login()
	case "register":
		register()
	case "refresh":
		refresh()
	case "addnote":
		addNote()
	case "delnote":
		delNote()
	case "updnote":
		updNote()
	case "addcard":
		addCard()
	case "addcred":
		addCred()
	case "getStream":
		getUpdateStream()
	case "addfile":
		addFile()
	case "delFile":
		delFile()
	}

}

func delFile() {
	fmt.Println("DEL FILE")

	ctx := metadata.AppendToOutgoingContext(context.Background(), "bearer", jwt)
	req := &proto.File{
		Id:       id,
		Serial:   serial,
		Name:     name,
		Filename: "asdasd",
		Metadata: []*proto.Meta{&meta1, &meta2},
	}
	_, err := docs.DeleteFile(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
}

func addFile() {
	fmt.Println("ADD FILE")
	file := "/tmp/test/big.dat"
	stat, _ := os.Stat(file)

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	chunk := &proto.FileChunk{
		Eof:  false,
		Data: nil,
	}

	chunkMessage := &proto.FileStream{ChunkedFile: &proto.FileStream_Chunk{Chunk: chunk}}
	fileMessage := &proto.FileStream{ChunkedFile: &proto.FileStream_File{File: &proto.File{
		Name:     "small",
		Metadata: []*proto.Meta{&meta1, &meta2},
		Filename: stat.Name(),
		Size:     stat.Size(),
	}}}
	ctx := metadata.AppendToOutgoingContext(context.Background(), "bearer", jwt)
	stream, err := docs.AddFile(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// send file meta
	err = stream.Send(fileMessage)
	if err != nil {
		log.Printf("Error on only send: %s", err.Error())
		return
	}
	bytes := make([]byte, 1<<18)
	for {
		n, err := f.Read(bytes)
		log.Printf("read %d", n)
		if err == io.EOF {
			chunk.Eof = true
			err = nil
		}
		if err != nil {
			log.Printf("failed to read file: %s", err.Error())
			return
		}
		chunk.Data = bytes[:n]
		err = stream.Send(chunkMessage)
		if err != nil {
			log.Printf("failed to send chunk: %s", err.Error())
			return
		}
		if chunk.Eof {
			_, err = stream.CloseAndRecv()
			if err != nil {
				log.Printf("failed to close stream: %s", err.Error())
				return
			}
			log.Println("sent!!")
			return

		}
	}
}

func getUpdateStream() {
	fmt.Println("UPDATE STREAM")
	ctx := metadata.AppendToOutgoingContext(context.Background(), "bearer", jwt)
	req := &proto.UpdateRequest{Serial: int64(serial)}
	stream, err := docs.GetUpdateStream(ctx, req)
	if err != nil {
		fmt.Println("error ", err)
		return
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("end of stream")
			return
		}
		if err != nil {
			fmt.Println("instream error: ", err)
			return
		}
		switch update := msg.Update.(type) {
		case *proto.UpdateResponse_Note:
			note := update.Note
			fmt.Println("note: ", note)
		case *proto.UpdateResponse_Credential:
			cred := update.Credential
			fmt.Println("cred: ", cred)
		case *proto.UpdateResponse_Card:
			card := update.Card
			fmt.Println("card: ", card)
		}
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
	_, err := auth.Register(context.Background(), req)
	if err != nil {
		panic(err)
	} else {
		log.Println("registered")
	}

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
	ctx := metadata.AppendToOutgoingContext(context.Background(), "bearer", jwt)
	req := &proto.Note{
		Name:     name,
		Metadata: []*proto.Meta{&meta1, &meta2},
		Text:     text,
	}
	_, err := docs.AddNote(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
}

func delNote() {
	fmt.Println("DEL NOTE")
	ctx := metadata.AppendToOutgoingContext(context.Background(), "bearer", jwt)
	req := &proto.Note{
		Id:       id,
		Serial:   serial,
		Name:     name,
		Metadata: []*proto.Meta{&meta1, &meta2},
		Text:     text,
	}
	_, err := docs.DeleteNote(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
}

func updNote() {
	fmt.Println("MOD NOTE")
	ctx := metadata.AppendToOutgoingContext(context.Background(), "bearer", jwt)
	req := &proto.Note{
		Id:       id,
		Serial:   serial,
		Name:     name,
		Metadata: []*proto.Meta{&meta1, &meta2},
		Text:     text,
	}
	_, err := docs.UpdateNote(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
}

func addCard() {
	fmt.Println("ADD CARD")
	ctx := metadata.AppendToOutgoingContext(context.Background(), "bearer", jwt)
	req := &proto.Card{
		Name:     name,
		Metadata: []*proto.Meta{&meta1, &meta2},
		Expires:  expires,
	}
	_, err := docs.AddCard(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
}

func addCred() {
	fmt.Println("ADD CRED")
	ctx := metadata.AppendToOutgoingContext(context.Background(), "bearer", jwt)
	req := &proto.Credential{
		Name:     name,
		Metadata: []*proto.Meta{&meta1, &meta2},
		Password: password,
	}
	_, err := docs.AddCredential(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
}
