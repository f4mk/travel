package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	apiDir      = "./api/specs"
	internalDir = "./internal/app/service"
)

func main() {
	modelDirs, err := os.ReadDir(apiDir)
	if err != nil {
		log.Fatalf("Failed to read directory %s: %v", apiDir, err)
	}
	for _, modelDir := range modelDirs {
		if modelDir.IsDir() {
			modelName := modelDir.Name()
			specPath := fmt.Sprintf("%s/%s/%s%s", apiDir, modelName, modelName, ".yaml")
			if _, err := os.Stat(specPath); !os.IsNotExist(err) {
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
