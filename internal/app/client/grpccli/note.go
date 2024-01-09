package grpccli

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"

	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/models"
)

// AddNote saves new Note on server
func (c *Client) AddNote(d models.Note) error {
	log.Println("grpc add note request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := proto.FromNote(d)
	if _, err := c.docs.AddNote(ctx, req); err != nil {
		log.Printf("grpc add note failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}

// UpdateNote updates Note on server
func (c *Client) UpdateNote(d models.Note) error {
	log.Println("grpc update note request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := proto.FromNote(d)
	if _, err := c.docs.UpdateNote(ctx, req); err != nil {
		log.Printf("grpc update note failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}

// DeleteNote deletes Note on server
func (c *Client) DeleteNote(d models.Note) error {
	log.Println("grpc delete note request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := proto.FromNote(d)
	if _, err := c.docs.DeleteNote(ctx, req); err != nil {
		log.Printf("grpc delete note failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}
