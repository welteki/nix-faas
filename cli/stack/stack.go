package stack

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ImageMetadata for function images build with nix
type ImageMetadata struct {
	Specifier string `yaml:"specifier,omitempty"`
	Source    string `yaml:"source,omitempty"`
}

// StackMetadata required by nix-faas
type StackMetadata struct {
	Images []ImageMetadata `yaml:"images,omitempty"`
}

// NixFaas configuration for stack.yaml
type NixFaas struct {
	StackMetadata StackMetadata `yaml:"x-nix-faas,omitempty"`
}

// ReadNixFaasConfig from the specified YAML file.
func ReadNixFaasConfig(yamlFile string) (NixFaas, error) {
	config := NixFaas{}

	configBytes, err := os.ReadFile(yamlFile)
	if err != nil {
		return config, fmt.Errorf("reading file %q: %w", yamlFile, err)
	}
	unmarshallErr := yaml.Unmarshal(configBytes, &config)
	if unmarshallErr != nil {
		return config, fmt.Errorf("reading configuration: %w", err)
	}
	return config, nil
}
