package types

import (
	"strings"
)

type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityMedium   Severity = "medium"
	SeverityLow      Severity = "low"
	SeverityWarning  Severity = "warning"
)

type ComplianceFramework string

const (
	ComplianceSOC2   ComplianceFramework = "soc2"
	CompliancePCIDSS ComplianceFramework = "pci-dss"
	ComplianceHIPAA  ComplianceFramework = "hipaa"
)

type Issue struct {
	ID          string              `json:"id"`
	RuleID      string              `json:"rule_id"`
	Severity    Severity            `json:"severity"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	File        string              `json:"file"`
	Line        int                 `json:"line,omitempty"`
	Column      int                 `json:"column,omitempty"`
	Resource    string              `json:"resource,omitempty"`
	Compliance  []ComplianceFramework `json:"compliance,omitempty"`
	Remediation Remediation         `json:"remediation"`
}

type Remediation struct {
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	CodeExample string   `json:"code_example,omitempty"`
}

type ScanResults struct {
	Issues []Issue `json:"issues"`
	Files  []string `json:"files_scanned"`
}

func NewScanResults() *ScanResults {
	return &ScanResults{
		Issues: make([]Issue, 0),
		Files:  make([]string, 0),
	}
}

func (sr *ScanResults) AddIssues(file string, issues []Issue) {
	sr.Files = append(sr.Files, file)
	sr.Issues = append(sr.Issues, issues...)
}

func (sr *ScanResults) FilterBySeverity(severity string) *ScanResults {
	filtered := NewScanResults()
	filtered.Files = sr.Files
	
	severityLower := strings.ToLower(severity)
	for _, issue := range sr.Issues {
		if strings.ToLower(string(issue.Severity)) == severityLower {
			filtered.Issues = append(filtered.Issues, issue)
		}
	}
	
	return filtered
}

func (sr *ScanResults) FilterByCompliance(framework string) *ScanResults {
	filtered := NewScanResults()
	filtered.Files = sr.Files
	
	frameworkLower := strings.ToLower(framework)
	for _, issue := range sr.Issues {
		for _, comp := range issue.Compliance {
			if strings.ToLower(string(comp)) == frameworkLower {
				filtered.Issues = append(filtered.Issues, issue)
				break
			}
		}
	}
	
	return filtered
}

type Summary struct {
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Warnings int `json:"warnings"`
	Total    int `json:"total"`
}

func (sr *ScanResults) Summary() Summary {
	summary := Summary{Total: len(sr.Issues)}
	
	for _, issue := range sr.Issues {
		switch issue.Severity {
		case SeverityCritical:
			summary.Critical++
		case SeverityHigh:
			summary.High++
		case SeverityMedium:
			summary.Medium++
		case SeverityLow:
			summary.Low++
		case SeverityWarning:
			summary.Warnings++
		}
	}
	
	return summary
}

// TerraformConfig represents parsed Terraform configuration
type TerraformConfig struct {
	Resources map[string]interface{} `json:"resources"`
	Variables map[string]interface{} `json:"variables"`
	Outputs   map[string]interface{} `json:"outputs"`
	Providers map[string]interface{} `json:"providers"`
}

// CloudFormationConfig represents parsed CloudFormation template
type CloudFormationConfig struct {
	Resources map[string]interface{} `json:"resources"`
	Parameters map[string]interface{} `json:"parameters"`
	Outputs   map[string]interface{} `json:"outputs"`
}

// KubernetesConfig represents parsed Kubernetes manifest
type KubernetesConfig struct {
	Kind     string                 `json:"kind"`
	Metadata map[string]interface{} `json:"metadata"`
	Spec     map[string]interface{} `json:"spec"`
}

