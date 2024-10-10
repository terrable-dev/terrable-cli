package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type TerrableToml struct {
	Environment map[string]string `toml:"environment"`
	Offline     OfflineConfig     `toml:"offline"`
}

type OfflineConfig struct {
	File   string `toml:"file"`
	Module string `toml:"module"`
	Port   string `toml:"port"`
}

func ParseTerrableToml() (*TerrableToml, error) {
	workingDir, _ := os.Getwd()

	// Attempt to find a .terrable.toml file in the active directory
	filePath := filepath.Join(workingDir, ".terrable.toml")
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
