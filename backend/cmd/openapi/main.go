package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gofiber/fiber/v2"

	"penguin-backend/internal/endpoints"
	"penguin-backend/internal/huma/fiberv2"
	"penguin-backend/internal/services"
)

func main() {
	container, err := services.CreateContainer()
	if err != nil {
		log.Fatalf("failed to initialize services: %v", err)
	}
	defer container.Cleanup()

	fiberApp := fiber.New()

	config := huma.DefaultConfig("Penguin ファイル情報管理API", "1.0.0")
	config.OpenAPI.Info.Description = "ファイルエントリの管理と閲覧のためのAPI"
	config.OpenAPI.Servers = []*huma.Server{
		{URL: "http://localhost:8080/api"},
		{URL: "https://localhost:8443/api"},
	}

	api := fiberv2.New(fiberApp, config)

	endpoints.SetupRoutes(fiberApp, api, *container)

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
