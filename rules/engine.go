package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iac-security-scanner/scanner/types"
	"gopkg.in/yaml.v3"
)

type Engine struct {
	rulesPath string
	rules     []Rule
}

type Rule struct {
	ID          string   `yaml:"id"`
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Severity    string   `yaml:"severity"`
	Type        string   `yaml:"type"` // terraform, cloudformation, kubernetes
	Compliance  []string `yaml:"compliance"`
	Conditions  []Condition `yaml:"conditions"`
	Remediation RemediationRule `yaml:"remediation"`
}

type Condition struct {
	Resource string            `yaml:"resource"`
	Property string            `yaml:"property"`
	Operator string            `yaml:"operator"` // equals, contains, exists, missing
	Value    interface{}       `yaml:"value"`
}

type RemediationRule struct {
	Description string   `yaml:"description"`
	Steps       []string `yaml:"steps"`
	CodeExample string   `yaml:"code_example"`
}

func NewEngine(rulesPath string) *Engine {
	return &Engine{
		rulesPath: rulesPath,
		rules:     make([]Rule, 0),
	}
}

func (e *Engine) LoadRules() error {
	if e.rulesPath == "" {
		e.rulesPath = "rules"
	}
	
	return filepath.Walk(e.rulesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() || !strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml") {
			return nil
		}
		
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read rule file %s: %w", path, err)
		}
		
		var rule Rule
		if err := yaml.Unmarshal(data, &rule); err != nil {
			return fmt.Errorf("failed to parse rule file %s: %w", path, err)
		}
		
		e.rules = append(e.rules, rule)
		return nil
	})
}

func (e *Engine) EvaluateTerraform(config *types.TerraformConfig, filePath string) []types.Issue {
	var issues []types.Issue
	
	for _, rule := range e.rules {
		if rule.Type != "terraform" {
			continue
		}
		
		if e.matchesRule(config, rule) {
			issue := e.createIssue(rule, filePath, "terraform")
			issues = append(issues, issue)
		}
	}
	
	return issues
}

func (e *Engine) EvaluateCloudFormation(config *types.CloudFormationConfig, filePath string) []types.Issue {
	var issues []types.Issue
	
	for _, rule := range e.rules {
		if rule.Type != "cloudformation" {
			continue
		}
		
		if e.matchesRule(config, rule) {
			issue := e.createIssue(rule, filePath, "cloudformation")
			issues = append(issues, issue)
		}
	}
	
	return issues
}

func (e *Engine) EvaluateKubernetes(config *types.KubernetesConfig, filePath string) []types.Issue {
	var issues []types.Issue
	
	for _, rule := range e.rules {
		if rule.Type != "kubernetes" {
			continue
		}
		
		if e.matchesRule(config, rule) {
			issue := e.createIssue(rule, filePath, "kubernetes")
			issues = append(issues, issue)
		}
	}
	
	return issues
}

func (e *Engine) matchesRule(config interface{}, rule Rule) bool {
	for _, condition := range rule.Conditions {
		if !e.evaluateCondition(config, condition) {
			return false
		}
	}
	return true
}

func (e *Engine) evaluateCondition(config interface{}, condition Condition) bool {
	// Simplified condition evaluation
	// In production, this would need more sophisticated path traversal
	switch config := config.(type) {
	case *types.TerraformConfig:
		return e.evaluateTerraformCondition(config, condition)
	case *types.CloudFormationConfig:
		return e.evaluateCloudFormationCondition(config, condition)
	case *types.KubernetesConfig:
		return e.evaluateKubernetesCondition(config, condition)
	}
	return false
}

func (e *Engine) evaluateTerraformCondition(config *types.TerraformConfig, condition Condition) bool {
	// Check if resource type exists (e.g., "aws_s3_bucket")
	if condition.Resource != "" {
		resourceExists := false
		for key := range config.Resources {
			// Key format: "resource_type.resource_name"
			if strings.HasPrefix(key, condition.Resource+".") {
				resourceExists = true
				break
			}
		}
		
		if condition.Operator == "exists" {
			return resourceExists
		}
		if condition.Operator == "missing" {
			return !resourceExists
		}
		
		// For equals/contains operators, check property values
		if condition.Property != "" && resourceExists {
			return e.checkTerraformProperty(config, condition)
		}
	}
	
	return false
}

func (e *Engine) checkTerraformProperty(config *types.TerraformConfig, condition Condition) bool {
	// Find the resource and check its property
	for key, resourceData := range config.Resources {
		if strings.HasPrefix(key, condition.Resource+".") {
			if resourceMap, ok := resourceData.(map[string]interface{}); ok {
				// Handle array properties (e.g., ingress rules)
				if arr, ok := resourceMap[strings.Split(condition.Property, ".")[0]].([]interface{}); ok {
					for _, item := range arr {
						if itemMap, ok := item.(map[string]interface{}); ok {
							// Check nested property in array item
							propPath := strings.Join(strings.Split(condition.Property, ".")[1:], ".")
							if propPath == "" {
								propPath = condition.Property
							}
							value := e.getNestedValue(itemMap, propPath)
							if e.compareValue(value, condition.Operator, condition.Value) {
								return true
							}
						}
					}
				} else {
					value := e.getNestedValue(resourceMap, condition.Property)
					if e.compareValue(value, condition.Operator, condition.Value) {
						return true
					}
				}
			}
		}
	}
	return false
}

func (e *Engine) getNestedValue(data map[string]interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	current := data
	
	for i, part := range parts {
		val, exists := current[part]
		if !exists {
			// Try to find in arrays
			for _, v := range current {
				if arr, ok := v.([]interface{}); ok {
					for _, item := range arr {
						if itemMap, ok := item.(map[string]interface{}); ok {
							if nestedVal, found := itemMap[part]; found {
								val = nestedVal
								exists = true
								if i == len(parts)-1 {
									return nestedVal
								}
								if nextMap, ok := nestedVal.(map[string]interface{}); ok {
									current = nextMap
									break
								}
							}
						}
					}
				}
			}
			if !exists {
				return nil
			}
		}
		
		if i == len(parts)-1 {
			return val
		}
		
		if nextMap, ok := val.(map[string]interface{}); ok {
			current = nextMap
		} else if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
			// Handle arrays - check first element or all elements
			if firstItem, ok := arr[0].(map[string]interface{}); ok {
				current = firstItem
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return nil
}

func (e *Engine) compareValue(actual interface{}, operator string, expected interface{}) bool {
	if actual == nil {
		return operator == "missing"
	}
	
	// Handle arrays
	if arr, ok := actual.([]interface{}); ok {
		for _, item := range arr {
			if e.compareValue(item, operator, expected) {
				return true
			}
		}
		return false
	}
	
	switch operator {
	case "equals":
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		return actualStr == expectedStr
	case "contains":
		actualStr := fmt.Sprintf("%v", actual)
		expectedStr := fmt.Sprintf("%v", expected)
		return strings.Contains(actualStr, expectedStr)
	case "exists":
		return actual != nil
	case "missing":
		return actual == nil
	}
	return false
}

func (e *Engine) evaluateCloudFormationCondition(config *types.CloudFormationConfig, condition Condition) bool {
	if condition.Resource != "" {
		// Check if any resource matches the type (e.g., AWS::EC2::Instance)
		resourceExists := false
		var resourceData map[string]interface{}
		
		for resourceName, resource := range config.Resources {
			if resourceMap, ok := resource.(map[string]interface{}); ok {
				if resourceType, ok := resourceMap["Type"].(string); ok {
					if resourceType == condition.Resource {
						resourceExists = true
						// Get Properties if they exist
						if props, ok := resourceMap["Properties"].(map[string]interface{}); ok {
							resourceData = props
						} else {
							resourceData = resourceMap
						}
						break
					}
				}
			}
		}
		
		if condition.Operator == "exists" {
			return resourceExists
		}
		if condition.Operator == "missing" {
			return !resourceExists
		}
		
		// For equals/contains, check properties
		if condition.Property != "" && resourceExists && resourceData != nil {
			value := e.getNestedValue(resourceData, condition.Property)
			return e.compareValue(value, condition.Operator, condition.Value)
		}
	}
	return false
}

func (e *Engine) evaluateKubernetesCondition(config *types.KubernetesConfig, condition Condition) bool {
	// Check based on Kubernetes resource type and properties
	if condition.Resource != "" && config.Kind != condition.Resource {
		// Also check for Deployment, StatefulSet, etc. that contain Pod specs
		if config.Kind == "Deployment" || config.Kind == "StatefulSet" || config.Kind == "DaemonSet" {
			// These resources have Pod templates, so we check them too
		} else {
			return false
		}
	}
	
	// Check properties in spec or metadata
	if condition.Property != "" {
		if condition.Operator == "missing" {
			return !e.hasKubernetesProperty(config, condition.Property)
		}
		if condition.Operator == "exists" {
			return e.hasKubernetesProperty(config, condition.Property)
		}
		
		// For equals/contains, get the value and compare
		value := e.getKubernetesValue(config, condition.Property)
		return e.compareValue(value, condition.Operator, condition.Value)
	}
	
	return true
}

func (e *Engine) getKubernetesValue(config *types.KubernetesConfig, property string) interface{} {
	parts := strings.Split(property, ".")
	current := make(map[string]interface{})
	
	if parts[0] == "spec" {
		current = config.Spec
		parts = parts[1:]
	} else if parts[0] == "metadata" {
		current = config.Metadata
		parts = parts[1:]
	}
	
	for i, part := range parts {
		val, exists := current[part]
		if !exists {
			// Check in arrays (e.g., containers)
			for _, v := range current {
				if arr, ok := v.([]interface{}); ok {
					for _, item := range arr {
						if itemMap, ok := item.(map[string]interface{}); ok {
							if nestedVal, found := itemMap[part]; found {
								val = nestedVal
								exists = true
								if i == len(parts)-1 {
									return nestedVal
								}
								if nextMap, ok := nestedVal.(map[string]interface{}); ok {
									current = nextMap
									break
								}
							}
						}
					}
				}
			}
			if !exists {
				return nil
			}
		}
		
		if i == len(parts)-1 {
			return val
		}
		
		if nextMap, ok := val.(map[string]interface{}); ok {
			current = nextMap
		} else if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
			// Handle arrays - check first element
			if firstItem, ok := arr[0].(map[string]interface{}); ok {
				current = firstItem
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return nil
}

func (e *Engine) hasKubernetesProperty(config *types.KubernetesConfig, property string) bool {
	return e.getKubernetesValue(config, property) != nil
}


func (e *Engine) createIssue(rule Rule, filePath string, resourceType string) types.Issue {
	severity := types.Severity(strings.ToLower(rule.Severity))
	if severity == "" {
		severity = types.SeverityMedium
	}
	
	complianceFrameworks := make([]types.ComplianceFramework, 0)
	for _, comp := range rule.Compliance {
		complianceFrameworks = append(complianceFrameworks, types.ComplianceFramework(strings.ToLower(comp)))
	}
	
	return types.Issue{
		ID:          fmt.Sprintf("%s-%d", rule.ID, len(complianceFrameworks)),
		RuleID:      rule.ID,
		Severity:    severity,
		Title:       rule.Name,
		Description: rule.Description,
		File:        filePath,
		Compliance:  complianceFrameworks,
		Remediation: types.Remediation{
			Description: rule.Remediation.Description,
			Steps:       rule.Remediation.Steps,
			CodeExample: rule.Remediation.CodeExample,
		},
	}
}

