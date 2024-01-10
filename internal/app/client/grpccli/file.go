package grpccli

import (
	"context"
	"errors"
	"io"
	"log"

	"google.golang.org/grpc/metadata"

	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/models"
)

// AddFile saves new File on server
func (c *Client) AddFile(d models.File, r io.Reader) error {
	log.Println("grpc add file request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	stream, err := c.docs.AddFile(ctx)
	if err != nil {
		return parseErr(err)
	}
	return c.sendFileToStream(d, r, stream)
}

// UpdateFile updates File on server with data
func (c *Client) UpdateFile(d models.File, r io.Reader) error {
	log.Println("grpc add file request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	stream, err := c.docs.UpdateFile(ctx)
	if err != nil {
		return parseErr(err)
	}
	return c.sendFileToStream(d, r, stream)
}

func (c *Client) sendFileToStream(d models.File, r io.Reader, stream proto.Docs_AddFileClient) error {
	chunk := &proto.FileChunk{}
	chunkMessage := &proto.FileStream{ChunkedFile: &proto.FileStream_Chunk{Chunk: chunk}}
	fileMessage := &proto.FileStream{ChunkedFile: &proto.FileStream_File{File: proto.FromFile(d)}}
	err := stream.Send(fileMessage)
	if err != nil {
		return parseErr(err)
	}
	chunkBytes := make([]byte, 1<<18)
	for {
		n, err := r.Read(chunkBytes)
		if err == io.EOF {
			chunk.Eof = true
			err = nil
		}
		if err != nil {
			return errors.New("fail read failed")
		}
		chunk.Data = chunkBytes[:n]
		err = stream.Send(chunkMessage)
		if err != nil {
			return parseErr(err)
		}
		if chunk.Eof {
			_, err = stream.CloseAndRecv()
			if err != nil {
				return parseErr(err)
			}
			return nil
		}
	}
}

// UpdateFileInfo updates FileInfo on server
func (c *Client) UpdateFileInfo(d models.File) error {
	log.Println("grpc update file request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	chunkMessage := &proto.FileStream{ChunkedFile: &proto.FileStream_Chunk{Chunk: &proto.FileChunk{Eof: true}}}
	fileMessage := &proto.FileStream{ChunkedFile: &proto.FileStream_File{File: proto.FromFile(d)}}
	stream, err := c.docs.UpdateFile(ctx)
	if err != nil {
		return parseErr(err)
	}
	if err := stream.Send(fileMessage); err != nil {
		log.Printf("grpc update file failed: %s", err.Error())
		return parseErr(err)
	}
	if err := stream.Send(chunkMessage); err != nil {
		log.Printf("grpc update file failed: %s", err.Error())
		return parseErr(err)
	}
	if _, err := stream.CloseAndRecv(); err != nil {
		log.Printf("grpc update file failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}

// DeleteFile deletes File on server
func (c *Client) DeleteFile(d models.File) error {
	log.Println("grpc delete file request")
	ctx, cancel := context.WithTimeout(context.Background(), c.dataTimeout)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "bearer", c.getToken())
	req := proto.FromFile(d)
	if _, err := c.docs.DeleteFile(ctx, req); err != nil {
		log.Printf("grpc delete file failed: %s", err.Error())
		return parseErr(err)
	}
	return nil
}
