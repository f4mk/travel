package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	apiDir      = "./api/specs"
	internalDir = "./internal/app/service"
)

func main() {
	modelDirs, err := os.ReadDir(apiDir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", apiDir, err)
	}
	for _, modelFile := range modelDirs {
		if !modelFile.IsDir() { // Ensure the entry is a file
			specPath := filepath.Join(apiDir, modelFile.Name())

			// Check if the file has a .yaml extension
			if filepath.Ext(specPath) == ".yaml" {
				modelName := strings.TrimSuffix(modelFile.Name(), ".yaml")
				outPath := fmt.Sprintf("%s/%s/model.go", internalDir, modelName)

				err := generateModel(specPath, modelName, outPath)
				if err != nil {
					log.Printf("Failed to generate model for spec %s: %v", specPath, err)
				}
			}
		}
	}

}

func generateModel(specPath, modelName, outPath string) error {
	fmt.Printf("oapi-codegen -package %s -generate types %s\n", modelName, specPath)
	cmd := exec.Command(
		"oapi-codegen",
		"-package", modelName,
		"-generate", "types",
		specPath)

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, out, 0644)
}
