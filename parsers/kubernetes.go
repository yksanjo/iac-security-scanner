package parsers

import (
	"fmt"
	"os"
	"strings"

	"github.com/iac-security-scanner/scanner/types"
	"gopkg.in/yaml.v3"
)

type KubernetesParser struct{}

func NewKubernetesParser() *KubernetesParser {
	return &KubernetesParser{}
}

func (p *KubernetesParser) Parse(path string) (*types.KubernetesConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	// Kubernetes files can contain multiple documents separated by ---
	documents := strings.Split(string(data), "---")
	
	// Parse the first document (or all if needed)
	var config types.KubernetesConfig
	
	for _, doc := range documents {
		doc = strings.TrimSpace(doc)
		if doc == "" {
			continue
		}
		
		var manifest map[string]interface{}
		if err := yaml.Unmarshal([]byte(doc), &manifest); err != nil {
			continue // Skip invalid documents
		}
		
		if kind, ok := manifest["kind"].(string); ok {
			config.Kind = kind
		}
		
		if metadata, ok := manifest["metadata"].(map[string]interface{}); ok {
			config.Metadata = metadata
		}
		
		if spec, ok := manifest["spec"].(map[string]interface{}); ok {
			config.Spec = spec
		}
		
		// For now, return the first valid document
		// In production, you might want to return all documents
		break
	}
	
	if config.Kind == "" {
		return nil, fmt.Errorf("invalid Kubernetes manifest: missing kind")
	}
	
	return &config, nil
}

