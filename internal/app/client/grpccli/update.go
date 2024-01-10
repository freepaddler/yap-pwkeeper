package grpccli

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc/metadata"

	"yap-pwkeeper/internal/pkg/grpc/proto"
)

func (c *Client) GetUpdateStream(serial int64, chData chan interface{}, chErr chan error) {
	defer func() {
		close(chData)
		close(chErr)
	}()
	log.Println("grpc update: started")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := &proto.UpdateRequest{Serial: serial}
	stream, err := c.docs.GetUpdateStream(ctx, req)
	if err != nil {
		log.Printf("grpc update: request failed: %s", err.Error())
		chErr <- parseErr(err)
		return
	}
	counter := 0
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Printf("grpc update: end of stream, received %d updates", counter)
			return
		}
		if err != nil {
			log.Printf("grpc update: stream error: %s", err.Error())
			chErr <- parseErr(err)
			return
		}
		switch update := msg.Update.(type) {
		case *proto.UpdateResponse_Note:
			note, err := update.Note.ToNote()
			if err != nil {
				log.Printf("grpc update: invalid note: %s", err)
				continue
			}
			//log.Printf("note %v", note)
			chData <- note
		case *proto.UpdateResponse_Credential:
			cred, err := update.Credential.ToCredential()
			if err != nil {
				log.Printf("grpc update: invalid credential: %s", err)
				continue
			}
			//log.Printf("credential %v", cred)
			chData <- cred
		case *proto.UpdateResponse_Card:
			card, err := update.Card.ToCard()
			if err != nil {
				log.Printf("grpc update: invalid card: %s", err)
				continue
			}
			//log.Printf("card %v", card)
			chData <- card
		case *proto.UpdateResponse_File:
			file, err := update.File.ToFile()
			if err != nil {
				log.Printf("grpc update: invalid file: %s", err)
				continue
			}
			log.Printf("file %v", file)
			chData <- file
		}
		counter++
	}
}
