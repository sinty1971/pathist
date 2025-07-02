package services

import (
	"fmt"
	"os"
	"path/filepath"
	"penguin-backend/internal/models"
	"strings"

	"gopkg.in/yaml.v3"
)

// CompanyService handles company data operations
type CompanyService struct {
	basePath string
}

// NewCompanyService creates a new CompanyService
func NewCompanyService() *CompanyService {
	// Default base path for company data
	homeDir, _ := os.UserHomeDir()
	basePath := filepath.Join(homeDir, "penguin", "豊田築炉", "1-会社情報")
	
	return &CompanyService{
		basePath: basePath,
	}
}

// SetBasePath sets the base path for company data
func (cs *CompanyService) SetBasePath(path string) {
	cs.basePath = path
}

// GetCompaniesFilePath returns the path to the companies data file
func (cs *CompanyService) GetCompaniesFilePath() string {
	return filepath.Join(cs.basePath, "companies.yaml")
}

// LoadCompanies loads all companies from the YAML file
func (cs *CompanyService) LoadCompanies() ([]models.Company, error) {
	filePath := cs.GetCompaniesFilePath()
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Return empty slice if file doesn't exist
		return []models.Company{}, nil
	}
	
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read companies file: %v", err)
	}
	
	var companies []models.Company
	if err := yaml.Unmarshal(data, &companies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal companies data: %v", err)
	}
	
	return companies, nil
}

// SaveCompanies saves all companies to the YAML file
func (cs *CompanyService) SaveCompanies(companies []models.Company) error {
	filePath := cs.GetCompaniesFilePath()
	
	// Ensure the directory exists
	if err := os.MkdirAll(cs.basePath, 0755); err != nil {
		return fmt.Errorf("failed to create companies directory: %v", err)
	}
	
	// Marshal companies to YAML
	data, err := yaml.Marshal(companies)
	if err != nil {
		return fmt.Errorf("failed to marshal companies data: %v", err)
	}
	
	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write companies file: %v", err)
	}
	
	return nil
}

// GetCompanyByID retrieves a company by its ID
func (cs *CompanyService) GetCompanyByID(id string) (*models.Company, error) {
	companies, err := cs.LoadCompanies()
	if err != nil {
		return nil, err
	}
	
	for _, company := range companies {
		if company.ID == id {
			return &company, nil
		}
	}
	
	return nil, fmt.Errorf("company with ID %s not found", id)
}

// GetCompanyByName retrieves a company by its name (case-insensitive)
func (cs *CompanyService) GetCompanyByName(name string) (*models.Company, error) {
	companies, err := cs.LoadCompanies()
	if err != nil {
		return nil, err
	}
	
	name = strings.TrimSpace(strings.ToLower(name))
	
	for _, company := range companies {
		if strings.ToLower(company.Name) == name || 
		   strings.ToLower(company.ShortName) == name ||
		   strings.ToLower(company.FullName) == name {
			return &company, nil
		}
	}
	
	return nil, nil // Return nil if not found (not an error)
}

// CreateCompany creates a new company
func (cs *CompanyService) CreateCompany(company models.Company) error {
	companies, err := cs.LoadCompanies()
	if err != nil {
		return err
	}
	
	// Check if company with same ID or name already exists
	for _, existing := range companies {
		if existing.ID == company.ID {
			return fmt.Errorf("company with ID %s already exists", company.ID)
		}
		if strings.EqualFold(existing.Name, company.Name) {
			return fmt.Errorf("company with name %s already exists", company.Name)
		}
	}
	
	// Add the new company
	companies = append(companies, company)
	
	return cs.SaveCompanies(companies)
}

// UpdateCompany updates an existing company
func (cs *CompanyService) UpdateCompany(company models.Company) error {
	companies, err := cs.LoadCompanies()
	if err != nil {
		return err
	}
	
	// Find and update the company
	found := false
	for i, existing := range companies {
		if existing.ID == company.ID {
			company.UpdateTimestamp()
			companies[i] = company
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("company with ID %s not found", company.ID)
	}
	
	return cs.SaveCompanies(companies)
}

// DeleteCompany deletes a company by ID
func (cs *CompanyService) DeleteCompany(id string) error {
	companies, err := cs.LoadCompanies()
	if err != nil {
		return err
	}
	
	// Find and remove the company
	found := false
	for i, company := range companies {
		if company.ID == id {
			companies = append(companies[:i], companies[i+1:]...)
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("company with ID %s not found", id)
	}
	
	return cs.SaveCompanies(companies)
}

// GetOrCreateCompanyByName gets an existing company by name or creates a new one
func (cs *CompanyService) GetOrCreateCompanyByName(name string) (*models.Company, error) {
	// Try to get existing company
	existing, err := cs.GetCompanyByName(name)
	if err != nil {
		return nil, err
	}
	
	if existing != nil {
		return existing, nil
	}
	
	// Create new company
	newCompany := models.NewCompany(name)
	if err := cs.CreateCompany(newCompany); err != nil {
		return nil, err
	}
	
	return &newCompany, nil
}

// ListCompanies returns all companies with optional filtering
func (cs *CompanyService) ListCompanies(filter string) ([]models.Company, error) {
	companies, err := cs.LoadCompanies()
	if err != nil {
		return nil, err
	}
	
	if filter == "" {
		return companies, nil
	}
	
	// Filter companies by name or business type
	var filtered []models.Company
	filter = strings.ToLower(filter)
	
	for _, company := range companies {
		if strings.Contains(strings.ToLower(company.Name), filter) ||
		   strings.Contains(strings.ToLower(company.FullName), filter) ||
		   strings.Contains(strings.ToLower(company.BusinessType), filter) {
			filtered = append(filtered, company)
		}
	}
	
	return filtered, nil
}

// GetCompanyStats returns statistics about companies
func (cs *CompanyService) GetCompanyStats() (map[string]interface{}, error) {
	companies, err := cs.LoadCompanies()
	if err != nil {
		return nil, err
	}
	
	stats := map[string]interface{}{
		"total_companies": len(companies),
		"business_types":  make(map[string]int),
		"complete_profiles": 0,
	}
	
	businessTypes := stats["business_types"].(map[string]int)
	completeProfiles := 0
	
	for _, company := range companies {
		// Count business types
		if company.BusinessType != "" {
			businessTypes[company.BusinessType]++
		}
		
		// Count complete profiles
		if company.IsComplete() {
			completeProfiles++
		}
	}
	
	stats["complete_profiles"] = completeProfiles
	stats["completion_rate"] = float64(completeProfiles) / float64(len(companies)) * 100
	
	return stats, nil
}