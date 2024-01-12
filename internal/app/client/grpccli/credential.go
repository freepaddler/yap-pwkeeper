package grpccli

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"

	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/models"
)

// AddCredential saves new Credential on server
func (c *Client) AddCredential(d models.Credential) error {
	log.Println("grpc add credential request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := proto.FromCredential(d)
	if _, err := c.docs.AddCredential(ctx, req); err != nil {
		log.Printf("grpc add credential failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}

// UpdateCredential updates Credential on server
func (c *Client) UpdateCredential(d models.Credential) error {
	log.Println("grpc update credential request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := proto.FromCredential(d)
	if _, err := c.docs.UpdateCredential(ctx, req); err != nil {
		log.Printf("grpc update credential failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}

// DeleteCredential deletes Credential on server
func (c *Client) DeleteCredential(d models.Credential) error {
	log.Println("grpc delete credential request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := proto.FromCredential(d)
	if _, err := c.docs.DeleteCredential(ctx, req); err != nil {
		log.Printf("grpc delete credential failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}
