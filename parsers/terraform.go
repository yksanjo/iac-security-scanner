package parsers

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/iac-security-scanner/scanner/types"
)

type TerraformParser struct {
	parser *hclparse.Parser
}

func NewTerraformParser() *TerraformParser {
	return &TerraformParser{
		parser: hclparse.NewParser(),
	}
}

func (p *TerraformParser) Parse(path string) (*types.TerraformConfig, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	file, diags := p.parser.ParseHCL(src, path)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse errors: %v", diags)
	}
	
	config := &types.TerraformConfig{
		Resources: make(map[string]interface{}),
		Variables: make(map[string]interface{}),
		Outputs:   make(map[string]interface{}),
		Providers: make(map[string]interface{}),
	}
	
	// Extract blocks
	body := file.Body
	attrs, _ := body.JustAttributes()
	
	// Parse blocks (simplified - in production, use proper HCL decoding)
	blocks := body.Blocks
	for _, block := range blocks {
		switch block.Type {
		case "resource":
			resourceType := block.Labels[0]
			resourceName := block.Labels[1]
			key := fmt.Sprintf("%s.%s", resourceType, resourceName)
			config.Resources[key] = extractBlockAttributes(block)
		case "variable":
			if len(block.Labels) > 0 {
				config.Variables[block.Labels[0]] = extractBlockAttributes(block)
			}
		case "output":
			if len(block.Labels) > 0 {
				config.Outputs[block.Labels[0]] = extractBlockAttributes(block)
			}
		case "provider":
			if len(block.Labels) > 0 {
				config.Providers[block.Labels[0]] = extractBlockAttributes(block)
			}
		}
	}
	
	_ = attrs // Use attributes if needed
	
	return config, nil
}

func extractBlockAttributes(block *hcl.Block) map[string]interface{} {
	attrs := make(map[string]interface{})
	
	body := block.Body
	blockAttrs, _ := body.JustAttributes()
	
	for name, attr := range blockAttrs {
		// Try to get string value, fallback to raw expression
		val, diags := attr.Expr.Value(&hcl.EvalContext{})
		if !diags.HasErrors() {
			attrs[name] = val.AsString()
		} else {
			// Fallback: store as raw expression string
			attrs[name] = attr.Expr.Range().String()
		}
	}
	
	// Also extract nested blocks
	for _, nestedBlock := range body.Blocks {
		key := nestedBlock.Type
		if len(nestedBlock.Labels) > 0 {
			key = fmt.Sprintf("%s.%s", key, nestedBlock.Labels[0])
		}
		attrs[key] = extractBlockAttributes(nestedBlock)
	}
	
	return attrs
}

