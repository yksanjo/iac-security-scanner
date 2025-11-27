package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iac-security-scanner/scanner/parsers"
	"github.com/iac-security-scanner/scanner/rules"
	"github.com/iac-security-scanner/scanner/types"
)

type Scanner struct {
	scanPath string
	rulesPath string
	ruleEngine *rules.Engine
}

func NewScanner(scanPath, rulesPath string) *Scanner {
	return &Scanner{
		scanPath: scanPath,
		rulesPath: rulesPath,
		ruleEngine: rules.NewEngine(rulesPath),
	}
}

func (s *Scanner) Scan() (*types.ScanResults, error) {
	results := types.NewScanResults()
	
	// Load rules
	if err := s.ruleEngine.LoadRules(); err != nil {
		return nil, fmt.Errorf("failed to load rules: %w", err)
	}
	
	// Walk directory and scan files
	err := filepath.Walk(s.scanPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		// Detect file type and parse
		var issues []types.Issue
		var parseErr error
		
		switch {
		case strings.HasSuffix(path, ".tf") || strings.HasSuffix(path, ".tf.json"):
			issues, parseErr = s.scanTerraform(path)
		case strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml"):
			// Try Kubernetes first, then CloudFormation
			if strings.Contains(path, "k8s") || strings.Contains(path, "kubernetes") {
				issues, parseErr = s.scanKubernetes(path)
			} else {
				issues, parseErr = s.scanCloudFormation(path)
			}
		case strings.HasSuffix(path, ".json"):
			// Could be CloudFormation or Terraform JSON
			if strings.Contains(strings.ToLower(path), "cloudformation") || 
			   strings.Contains(strings.ToLower(path), "template") {
				issues, parseErr = s.scanCloudFormation(path)
			} else if strings.Contains(path, ".tf.json") {
				issues, parseErr = s.scanTerraform(path)
			}
		}
		
		if parseErr != nil {
			// Log but don't fail the entire scan
			fmt.Printf("⚠️  Warning: Failed to parse %s: %v\n", path, parseErr)
			return nil
		}
		
		results.AddIssues(path, issues)
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}
	
	return results, nil
}

func (s *Scanner) scanTerraform(path string) ([]types.Issue, error) {
	parser := parsers.NewTerraformParser()
	config, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}
	
	return s.ruleEngine.EvaluateTerraform(config, path), nil
}

func (s *Scanner) scanCloudFormation(path string) ([]types.Issue, error) {
	parser := parsers.NewCloudFormationParser()
	config, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}
	
	return s.ruleEngine.EvaluateCloudFormation(config, path), nil
}

func (s *Scanner) scanKubernetes(path string) ([]types.Issue, error) {
	parser := parsers.NewKubernetesParser()
	config, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}
	
	return s.ruleEngine.EvaluateKubernetes(config, path), nil
}

