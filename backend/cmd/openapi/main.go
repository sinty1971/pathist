package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"

	"penguin-backend/internal/app"
	"penguin-backend/internal/huma/fiberv2"
	"penguin-backend/internal/routes"
	svc "penguin-backend/internal/services"
)

func main() {
	serviceOptions := app.DefaultServiceOptions
	if dataRoot := os.Getenv("PENGUIN_DATA_ROOT"); dataRoot != "" {
		serviceOptions.FileFolderPath = dataRoot
		serviceOptions.CompanyFolderPath = filepath.Join(dataRoot, "豊田築炉", "1 会社")
		serviceOptions.KojiFolderPath = filepath.Join(dataRoot, "豊田築炉", "2 工事")
	}
	rootService := svc.CreateRootService()

	fileService := &svc.FileService{}
	if err := rootService.AddService(fileService, svc.WithPath(serviceOptions.FileFolderPath)); err != nil {
		log.Fatalf("failed to initialize FileService: %v", err)
	}

	companyService := &svc.CompanyService{}
	if err := rootService.AddService(
		companyService,
		svc.WithPath(serviceOptions.CompanyFolderPath),
		svc.WithFileName(serviceOptions.DatabaseFilename),
	); err != nil {
		log.Fatalf("failed to initialize CompanyService: %v", err)
	}

	kojiService := &svc.KojiService{}
	if err := rootService.AddService(
		kojiService,
		svc.WithPath(serviceOptions.KojiFolderPath),
		svc.WithFileName(serviceOptions.DatabaseFilename),
	); err != nil {
		log.Fatalf("failed to initialize KojiService: %v", err)
	}

	services := app.ServiceContainer{
		Root:    rootService,
		File:    fileService,
		Company: companyService,
		Koji:    kojiService,
	}
	defer services.Root.Cleanup()

	fiberApp := fiber.New()

	config := huma.DefaultConfig("Penguin ファイル情報管理API", "1.0.0")
	config.OpenAPI.Info.Description = "ファイルエントリの管理と閲覧のためのAPI"
	config.OpenAPI.Servers = []*huma.Server{
		{URL: "http://localhost:8080/api"},
		{URL: "https://localhost:8443/api"},
	}

	api := fiberv2.New(fiberApp, config)

	routes.SetupRoutes(fiberApp, api, services)

	spec := api.OpenAPI()

	outputDir := filepath.Join("..", "schemas")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("failed to create schemas directory: %v", err)
	}

	writeJSON(filepath.Join(outputDir, "openapi.json"), spec)
	writeYAML(filepath.Join(outputDir, "openapi.yaml"), spec)

	log.Println("OpenAPI 3.1 specification generated (JSON & YAML)")
}

func writeJSON(path string, spec *huma.OpenAPI) {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal OpenAPI JSON: %v", err)
	}
	writeBytes(path, data)
}

func writeYAML(path string, spec *huma.OpenAPI) {
	data, err := spec.YAML()
	if err != nil {
		log.Fatalf("failed to marshal OpenAPI YAML: %v", err)
	}
	writeBytes(path, data)
}

func writeBytes(path string, data []byte) {
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Fatalf("failed to write %s: %v", path, err)
	}
	log.Printf("wrote %s", path)
}
