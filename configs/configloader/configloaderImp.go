package configloader

import (
	"os"

	"gopkg.in/yaml.v2"
)

type YAMLConfigLoader struct {
	FilePath string
}

func NewYAMLConfigLoader(filePath string) *YAMLConfigLoader {
	return &YAMLConfigLoader{FilePath: filePath}
}

func (loader *YAMLConfigLoader) LoadConfig() (*Config, error) {
	data, err := os.ReadFile(loader.FilePath)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
