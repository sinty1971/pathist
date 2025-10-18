package rpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	penguinv1 "penguin-backend/gen/penguin/v1"
	penguinv1connect "penguin-backend/gen/penguin/v1/penguinv1connect"
	"penguin-backend/internal/models"
	"penguin-backend/internal/services"
)

// KojiServiceHandler bridges existing KojiService logic to Connect handlers.
type KojiServiceHandler struct {
	penguinv1connect.UnimplementedKojiServiceHandler
	kojiService *services.KojiService
}

func NewKojiServiceHandler(service *services.KojiService) *KojiServiceHandler {
	return &KojiServiceHandler{
		kojiService: service,
	}
}

func (h *KojiServiceHandler) ListKojies(ctx context.Context, req *connect.Request[penguinv1.ListKojiesRequest]) (*connect.Response[penguinv1.ListKojiesResponse], error) {
	_ = req // 現状フィルター未対応

	kojies := h.kojiService.GetKojies()
	items := make([]*penguinv1.Koji, 0, len(kojies))
	for i := range kojies {
		k := kojies[i]
		items = append(items, convertModelKoji(&k))
	}

	res := penguinv1.ListKojiesResponse_builder{
		Kojies: items,
	}.Build()

	return connect.NewResponse(res), nil
}

func (h *KojiServiceHandler) GetKoji(ctx context.Context, req *connect.Request[penguinv1.GetKojiRequest]) (*connect.Response[penguinv1.GetKojiResponse], error) {
	folderName := req.Msg.GetFolderName()
	if folderName == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("folder_name is required"))
	}

	koji, err := h.kojiService.GetKoji(folderName)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := penguinv1.GetKojiResponse_builder{
		Koji: convertModelKoji(koji),
	}.Build()

	return connect.NewResponse(res), nil
}

func (h *KojiServiceHandler) UpdateKoji(ctx context.Context, req *connect.Request[penguinv1.UpdateKojiRequest]) (*connect.Response[penguinv1.UpdateKojiResponse], error) {
	if req.Msg.GetKoji() == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("koji is required"))
	}

	modelKoji := convertProtoKoji(req.Msg.GetKoji())
	if err := h.kojiService.Update(modelKoji); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := penguinv1.UpdateKojiResponse_builder{
		Koji: convertModelKoji(modelKoji),
	}.Build()

	return connect.NewResponse(res), nil
}

func (h *KojiServiceHandler) UpdateKojiStandardFiles(ctx context.Context, req *connect.Request[penguinv1.UpdateKojiStandardFilesRequest]) (*connect.Response[penguinv1.UpdateKojiStandardFilesResponse], error) {
	if req.Msg.GetKoji() == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("koji is required"))
	}

	modelKoji := convertProtoKoji(req.Msg.GetKoji())
	folderName := modelKoji.GetFolderName()
	if folderName == "" {
		res := penguinv1.UpdateKojiStandardFilesResponse_builder{}.Build()
		return connect.NewResponse(res), nil
	}

	updated, err := h.kojiService.GetKoji(folderName)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := penguinv1.UpdateKojiStandardFilesResponse_builder{
		Koji: convertModelKoji(updated),
	}.Build()

	return connect.NewResponse(res), nil
}

func convertModelKoji(src *models.Koji) *penguinv1.Koji {
	if src == nil {
		return nil
	}

	files := make([]*penguinv1.FileInfo, 0, len(src.RequiredFiles))
	for i := range src.RequiredFiles {
		f := src.RequiredFiles[i]
		files = append(files, convertModelFileInfo(&f))
	}

	return penguinv1.Koji_builder{
		Id:            src.ID,
		Status:        src.Status,
		TargetFolder:  src.TargetFolder,
		StartDate:     toProtoTimestamp(src.StartDate),
		CompanyName:   src.CompanyName,
		LocationName:  src.LocationName,
		EndDate:       toProtoTimestamp(src.EndDate),
		Description:   src.Description,
		Tags:          append([]string(nil), src.Tags...),
		RequiredFiles: files,
	}.Build()
}

func convertProtoKoji(src *penguinv1.Koji) *models.Koji {
	if src == nil {
		return nil
	}

	model := &models.Koji{
		ID:           src.GetId(),
		Status:       src.GetStatus(),
		TargetFolder: src.GetTargetFolder(),
		StartDate:    toModelTimestamp(src.GetStartDate()),
		CompanyName:  src.GetCompanyName(),
		LocationName: src.GetLocationName(),
		EndDate:      toModelTimestamp(src.GetEndDate()),
		Description:  src.GetDescription(),
		Tags:         append([]string(nil), src.GetTags()...),
	}

	required := src.GetRequiredFiles()
	if len(required) > 0 {
		model.RequiredFiles = make([]models.FileInfo, 0, len(required))
		for _, item := range required {
			model.RequiredFiles = append(model.RequiredFiles, convertProtoFileInfo(item))
		}
	}

	return model
}
