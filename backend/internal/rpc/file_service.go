package rpc

import (
	"context"

	"connectrpc.com/connect"
	penguinv1 "penguin-backend/gen/penguin/v1"
	penguinv1connect "penguin-backend/gen/penguin/v1/penguinv1connect"
	"penguin-backend/internal/models"
	"penguin-backend/internal/services"
)

// FileServiceHandler exposes FileService operations via Connect handlers.
type FileServiceHandler struct {
	penguinv1connect.UnimplementedFileServiceHandler
	fileService *services.FileService
}

func NewFileServiceHandler(service *services.FileService) *FileServiceHandler {
	return &FileServiceHandler{
		fileService: service,
	}
}

func (h *FileServiceHandler) ListFiles(ctx context.Context, req *connect.Request[penguinv1.ListFilesRequest]) (*connect.Response[penguinv1.ListFilesResponse], error) {
	_ = ctx

	var (
		entries []models.FileInfo
		err     error
	)

	path := req.Msg.GetPath()
	if path == "" {
		entries, err = h.fileService.GetFileInfos()
	} else {
		entries, err = h.fileService.GetFileInfos(path)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	files := make([]*penguinv1.FileInfo, 0, len(entries))
	for i := range entries {
		files = append(files, convertModelFileInfo(&entries[i]))
	}

	res := penguinv1.ListFilesResponse_builder{
		Files: files,
	}.Build()

	return connect.NewResponse(res), nil
}

func (h *FileServiceHandler) GetFileBasePath(ctx context.Context, _ *connect.Request[penguinv1.GetFileBasePathRequest]) (*connect.Response[penguinv1.GetFileBasePathResponse], error) {
	_ = ctx

	res := penguinv1.GetFileBasePathResponse_builder{
		BasePath: h.fileService.TargetFolder,
	}.Build()

	return connect.NewResponse(res), nil
}
