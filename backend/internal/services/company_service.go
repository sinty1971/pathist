package services

import (
	"context"
	"errors"

	penguinv1 "penguin-backend/gen/penguin/v1"
	penguinv1connect "penguin-backend/gen/penguin/v1/penguinv1connect"
	"penguin-backend/internal/models"
	"penguin-backend/internal/services"

	"connectrpc.com/connect"
)

// CompanyServiceHandler bridges CompanyService logic to Connect handlers.
type CompanyServiceHandler struct {
	penguinv1connect.UnimplementedCompanyServiceHandler
	companyService *services.CompanyService
}

func NewCompanyServiceHandler(service *services.CompanyService) *CompanyServiceHandler {
	return &CompanyServiceHandler{
		companyService: service,
	}
}

func (h *CompanyServiceHandler) ListCompanies(ctx context.Context, _ *connect.Request[penguinv1.ListCompaniesRequest]) (*connect.Response[penguinv1.ListCompaniesResponse], error) {
	companies := h.companyService.GetCompanies()
	items := make([]*penguinv1.Company, 0, len(companies))
	for i := range companies {
		company := companies[i]
		items = append(items, convertModelCompany(&company))
	}

	res := penguinv1.ListCompaniesResponse_builder{
		Companies: items,
	}.Build()

	return connect.NewResponse(res), nil
}

func (h *CompanyServiceHandler) GetCompany(ctx context.Context, req *connect.Request[penguinv1.GetCompanyRequest]) (*connect.Response[penguinv1.GetCompanyResponse], error) {
	id := req.Msg.GetId()
	if id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}

	company, err := h.companyService.GetCompanyByID(id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	res := penguinv1.GetCompanyResponse_builder{
		Company: convertModelCompany(company),
	}.Build()

	return connect.NewResponse(res), nil
}

func (h *CompanyServiceHandler) UpdateCompany(ctx context.Context, req *connect.Request[penguinv1.UpdateCompanyRequest]) (*connect.Response[penguinv1.UpdateCompanyResponse], error) {
	if req.Msg.GetCompany() == nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("company is required"))
	}

	modelCompany := convertProtoCompany(req.Msg.GetCompany())
	if modelCompany.TargetFolder == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("target_folder is required"))
	}

	if err := h.companyService.Update(modelCompany); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	reloaded, err := h.companyService.GetCompany(modelCompany.TargetFolder)
	if err == nil {
		modelCompany = reloaded
	}

	res := penguinv1.UpdateCompanyResponse_builder{
		Company: convertModelCompany(modelCompany),
	}.Build()

	return connect.NewResponse(res), nil
}

func (h *CompanyServiceHandler) ListCompanyCategories(ctx context.Context, _ *connect.Request[penguinv1.ListCompanyCategoriesRequest]) (*connect.Response[penguinv1.ListCompanyCategoriesResponse], error) {
	categories := h.companyService.Categories()
	items := make([]*penguinv1.CompanyCategoryInfo, 0, len(categories))
	for _, category := range categories {
		items = append(items, convertModelCompanyCategory(category))
	}

	res := penguinv1.ListCompanyCategoriesResponse_builder{
		Categories: items,
	}.Build()

	return connect.NewResponse(res), nil
}

func convertModelCompany(src *models.Company) *penguinv1.Company {
	if src == nil {
		return nil
	}

	return src.Proto()
}

func convertProtoCompany(src *penguinv1.Company) *models.Company {
	if src == nil {
		return nil
	}

	return models.NewCompanyFromProto(src)
}

func convertModelCompanyCategory(src models.CompanyCategory) *penguinv1.CompanyCategoryInfo {
	return penguinv1.CompanyCategoryInfo_builder{
		Code:  string(src.Index),
		Label: src.Label,
	}.Build()
}
