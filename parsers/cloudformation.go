package parsers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iac-security-scanner/scanner/types"
	"gopkg.in/yaml.v3"
)

type CloudFormationParser struct{}

func NewCloudFormationParser() *CloudFormationParser {
	return &CloudFormationParser{}
}

func (p *CloudFormationParser) Parse(path string) (*types.CloudFormationConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	var template map[string]interface{}
	
	// Try YAML first
	if err := yaml.Unmarshal(data, &template); err != nil {
		// Fall back to JSON
		if err := json.Unmarshal(data, &template); err != nil {
			return nil, fmt.Errorf("failed to parse CloudFormation template: %w", err)
		}
	}
	
	config := &types.CloudFormationConfig{
		Resources: make(map[string]interface{}),
		Parameters: make(map[string]interface{}),
		Outputs:   make(map[string]interface{}),
	}
	
	if resources, ok := template["Resources"].(map[string]interface{}); ok {
		config.Resources = resources
	}
	
	if parameters, ok := template["Parameters"].(map[string]interface{}); ok {
		config.Parameters = parameters
	}
	
	if outputs, ok := template["Outputs"].(map[string]interface{}); ok {
		config.Outputs = outputs
	}
	
	return config, nil
}

