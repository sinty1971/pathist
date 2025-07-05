package handlers

import (
	"penguin-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

type CompanyHandler struct {
	CompanyService *services.CompanyService
}

func NewCompanyHandler(companyService *services.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		CompanyService: companyService,
	}
}

func (h *CompanyHandler) GetCompany(c *fiber.Ctx) error {
	company, err := h.CompanyService.GetCompanyByID(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get company",
			"message": err.Error(),
		})
	}
	return c.JSON(company)
}

func (h *CompanyHandler) GetCompanies(c *fiber.Ctx) error {
	companies := h.CompanyService.GetCompanies()
	return c.JSON(companies)
}
