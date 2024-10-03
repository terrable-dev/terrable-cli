package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type TerrableToml struct {
	Environment map[string]interface{} `toml:"environment"`
	Offline     OfflineConfig          `toml:"offline"`
}

type OfflineConfig struct {
	File   string `toml:"file"`
	Module string `toml:"module"`
}

func ParseTerrableToml() (*TerrableToml, error) {
	// Attempt to find a .terrable.toml file in the active directory
	workingDirectory, _ := os.Getwd()

	filePath := filepath.Join(workingDirectory, ".terrable.toml")
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return &TerrableToml{}, nil
	}

	// Read and parse terrable.toml
	var config TerrableToml

	_, err = toml.DecodeFile(filePath, &config)

	if err != nil {
		return nil, fmt.Errorf("failed to parse .terrable.toml file: %w", err)
	}

	return &config, nil
}
