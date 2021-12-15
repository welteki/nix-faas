package stack

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
