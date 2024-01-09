package grpcapi

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"yap-pwkeeper/internal/pkg/grpc/proto"
	"yap-pwkeeper/internal/pkg/logger"
	"yap-pwkeeper/internal/pkg/models"
)

func (w DocsHandlers) AddFile(stream proto.Docs_AddFileServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("add file stream request")
	return w.withFileFromStream(stream, w.docs.AddFile)
}

func (w DocsHandlers) UpdateFile(stream proto.Docs_UpdateFileServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("update file stream request")
	return w.withFileFromStream(stream, w.docs.AddFile)
}

// withFileFromStream gets file from stream and executes passed function
func (w DocsHandlers) withFileFromStream(stream proto.Docs_AddFileServer, fn func(ctx context.Context, file models.File) error) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("process file from stream")
	var file models.File
	eof := true
	hash := sha256.New()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Debug("file stream ended")
			// no chunk with oef flag
			if !eof {
				return status.Error(codes.InvalidArgument, "incomplete file data")
			}
			// processing result
			file.Size = int64(len(file.Data))
			file.Sha265 = fmt.Sprintf("%x", hash.Sum(nil))
			file.UserId, _ = logger.GetUserId(ctx)
			if err := fn(ctx, file); err != nil {
				return respErr(ctx, err)
			}
			return stream.SendAndClose(&proto.Empty{})
			//return file, err
		}
		if err != nil {
			return err
		}
		switch req.ChunkedFile.(type) {
		case *proto.FileStream_File:
			log.Debug("file message")
			if !eof {
				return status.Error(codes.InvalidArgument, "unexpected end of file")
			}
			file, err = req.ChunkedFile.(*proto.FileStream_File).File.ToFile()
			if err != nil {
				return status.Error(codes.InvalidArgument, err.Error())
			}
		case *proto.FileStream_Chunk:
			log.Debug("chunk message")
			if file.Filename == "" {
				return status.Error(codes.InvalidArgument, "chunk sent before file")
			}
			chunk := req.ChunkedFile.(*proto.FileStream_Chunk).Chunk
			eof = chunk.Eof
			file.Data = append(file.Data, chunk.Data...)
			if len(file.Data) > 14<<20 {
				return status.Error(codes.InvalidArgument, "file too big")
			}
			_, err = hash.Write(chunk.Data)
			if err != nil {
				logger.Log().WithErr(err).Error("failed to get file hash")
			}
		}
	}
}

func (w DocsHandlers) DeleteFile(ctx context.Context, in *proto.File) (*proto.Empty, error) {
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("delete file request")
	file, err := in.ToFile()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}
	file.UserId, _ = logger.GetUserId(ctx)
	if err := w.docs.DeleteFile(ctx, file); err != nil {
		return nil, respErr(ctx, err)
	}
	return &proto.Empty{}, nil
}

func (w DocsHandlers) GetFile(in *proto.DocumentRequest, stream proto.Docs_GetFileServer) error {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()
	log := logger.Log().WithCtxRequestId(ctx).WithCtxUserId(ctx)
	log.Debug("update stream request")
	userId, _ := logger.GetUserId(ctx)
	file, err := w.docs.GetFile(ctx, in.GetId(), userId)
	if err != nil {
		return respErr(ctx, err)
	}
	chunk := &proto.FileChunk{}
	chunkMessage := &proto.FileStream{ChunkedFile: &proto.FileStream_Chunk{Chunk: chunk}}
	fileMessage := &proto.FileStream{ChunkedFile: &proto.FileStream_File{File: proto.FromFile(file)}}
	err = stream.Send(fileMessage)
	if err != nil {
		return err
	}
	r := bytes.NewReader(file.Data)
	chunkBytes := make([]byte, 1<<18)
	for {
		n, err := r.Read(chunkBytes)
		if err == io.EOF {
			chunk.Eof = true
			err = nil
		}
		if err != nil {
			log.WithErr(err).Error("failed to read file bytes")
			return status.Error(codes.Internal, "")
		}
		chunk.Data = chunkBytes[:n]
		err = stream.Send(chunkMessage)
		if err != nil {
			return err
		}
		if chunk.Eof {
			return nil
		}
	}
}
