package grpccli

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"

	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/models"
)

// AddCard saves new Card on server
func (c *Client) AddCard(d models.Card) error {
	log.Println("grpc add card request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := proto.FromCard(d)
	if _, err := c.docs.AddCard(ctx, req); err != nil {
		log.Printf("grpc add card failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}

// UpdateCard updates card on server
func (c *Client) UpdateCard(d models.Card) error {
	log.Println("grpc update card request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := proto.FromCard(d)
	if _, err := c.docs.UpdateCard(ctx, req); err != nil {
		log.Printf("grpc update card failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}

// DeleteCard deletes card on server
func (c *Client) DeleteCard(d models.Card) error {
	log.Println("grpc delete card request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := proto.FromCard(d)
	if _, err := c.docs.DeleteCard(ctx, req); err != nil {
		log.Printf("grpc delete card failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}
