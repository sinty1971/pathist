package models

import (
	"fmt"
	"strings"
	"time"
)

// Company represents a construction company with detailed information
// @Description Construction company information with contact details and business info
type Company struct {
	// Basic company information
	ID   string `json:"id" yaml:"id" example:"TC001"`
	Name string `json:"name" yaml:"name" example:"豊田築炉"`
	
	// Company details
	FullName        string `json:"full_name,omitempty" yaml:"full_name" example:"豊田築炉株式会社"`
	ShortName       string `json:"short_name,omitempty" yaml:"short_name" example:"豊田築炉"`
	CompanyType     string `json:"company_type,omitempty" yaml:"company_type" example:"株式会社"`
	RegistrationNum string `json:"registration_num,omitempty" yaml:"registration_num" example:"1234567890"`
	
	// Contact information
	PostalCode string `json:"postal_code,omitempty" yaml:"postal_code" example:"456-0001"`
	Address    string `json:"address,omitempty" yaml:"address" example:"愛知県名古屋市熱田区三本松町1-1"`
	Phone      string `json:"phone,omitempty" yaml:"phone" example:"052-681-8111"`
	Fax        string `json:"fax,omitempty" yaml:"fax" example:"052-681-8112"`
	Email      string `json:"email,omitempty" yaml:"email" example:"info@toyoda-ro.co.jp"`
	Website    string `json:"website,omitempty" yaml:"website" example:"https://www.toyoda-ro.co.jp"`
	
	// Business information
	BusinessType string   `json:"business_type,omitempty" yaml:"business_type" example:"工業炉製造"`
	Services     []string `json:"services,omitempty" yaml:"services" example:"['工業炉設計', '工業炉製造', '工業炉メンテナンス']"`
	Specialties  []string `json:"specialties,omitempty" yaml:"specialties" example:"['高温炉', '熱処理炉', '焼成炉']"`
	
	// Timestamps
	CreatedAt Timestamp `json:"created_at,omitempty" yaml:"created_at"`
	UpdatedAt Timestamp `json:"updated_at,omitempty" yaml:"updated_at"`
	
	// Additional metadata
	Notes  string            `json:"notes,omitempty" yaml:"notes" example:"主要取引先企業"`
	Tags   []string          `json:"tags,omitempty" yaml:"tags" example:"['元請け', '製造業', '愛知県']"`
	Custom map[string]string `json:"custom,omitempty" yaml:"custom" example:"{'担当者': '田中太郎', '契約種別': '年間契約'}"`
	
	// Project relationship
	ProjectCount int `json:"project_count,omitempty" yaml:"-"` // Number of associated projects
}

// NewCompany creates a new Company from basic information
func NewCompany(name string) Company {
	now := Timestamp{Time: time.Now()}
	
	// Generate company ID from name
	id := GenerateCompanyID(name)
	
	company := Company{
		ID:           id,
		Name:         name,
		ShortName:    name,
		BusinessType: "建設業",
		CreatedAt:    now,
		UpdatedAt:    now,
		Tags:         []string{"会社", "建設"},
		Custom:       make(map[string]string),
	}
	
	return company
}

// GenerateCompanyID generates a unique ID for the company
func GenerateCompanyID(name string) string {
	// Remove spaces and generate ID from company name
	cleanName := strings.ReplaceAll(name, " ", "")
	idSource := fmt.Sprintf("company_%s_%d", cleanName, time.Now().Unix())
	id := NewIDFromString(idSource)
	return fmt.Sprintf("C%s", id.Len5())
}

// UpdateTimestamp updates the UpdatedAt timestamp to current time
func (c *Company) UpdateTimestamp() {
	c.UpdatedAt = Timestamp{Time: time.Now()}
}

// AddService adds a service to the company's services list
func (c *Company) AddService(service string) {
	for _, existing := range c.Services {
		if existing == service {
			return // Service already exists
		}
	}
	c.Services = append(c.Services, service)
	c.UpdateTimestamp()
}

// AddSpecialty adds a specialty to the company's specialties list
func (c *Company) AddSpecialty(specialty string) {
	for _, existing := range c.Specialties {
		if existing == specialty {
			return // Specialty already exists
		}
	}
	c.Specialties = append(c.Specialties, specialty)
	c.UpdateTimestamp()
}

// AddTag adds a tag to the company's tags list
func (c *Company) AddTag(tag string) {
	for _, existing := range c.Tags {
		if existing == tag {
			return // Tag already exists
		}
	}
	c.Tags = append(c.Tags, tag)
	c.UpdateTimestamp()
}

// SetCustomField sets a custom field value
func (c *Company) SetCustomField(key, value string) {
	if c.Custom == nil {
		c.Custom = make(map[string]string)
	}
	c.Custom[key] = value
	c.UpdateTimestamp()
}

// GetCustomField gets a custom field value
func (c *Company) GetCustomField(key string) (string, bool) {
	if c.Custom == nil {
		return "", false
	}
	value, exists := c.Custom[key]
	return value, exists
}

// IsComplete checks if the company has all required information
func (c *Company) IsComplete() bool {
	return c.Name != "" && c.Address != "" && c.Phone != ""
}

// GetDisplayName returns the display name for the company (FullName if available, otherwise Name)
func (c *Company) GetDisplayName() string {
	if c.FullName != "" {
		return c.FullName
	}
	return c.Name
}

// GetShortDisplayName returns a short display name for the company
func (c *Company) GetShortDisplayName() string {
	if c.ShortName != "" {
		return c.ShortName
	}
	return c.Name
}