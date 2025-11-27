package cmd

import (
	"fmt"

	"github.com/iac-security-scanner/scanner/reports"
	"github.com/iac-security-scanner/scanner/scanner"
	"github.com/spf13/cobra"
)

var (
	complianceFramework string
	outputFormat       string
	outputPath         string
	rulesPath          string
	severityFilter     string
)

var rootCmd = &cobra.Command{
	Use:   "iac-audit",
	Short: "Infrastructure-as-Code Security Scanner",
	Long: `A comprehensive security auditing tool for Infrastructure-as-Code files.
Scans Terraform, CloudFormation, and Kubernetes YAML for security misconfigurations
and generates compliance reports.`,
	Version: "1.0.0",
}

var scanCmd = &cobra.Command{
	Use:   "scan [path]",
	Short: "Scan IaC files for security issues",
	Long: `Scan Infrastructure-as-Code files in the specified directory
for security misconfigurations and compliance violations.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scanPath := args[0]
		
		fmt.Printf("🔍 Scanning: %s\n", scanPath)
		fmt.Printf("📋 Compliance: %s\n", complianceFramework)
		fmt.Printf("📄 Output: %s (%s)\n\n", outputPath, outputFormat)
		
		// Initialize scanner
		sc := scanner.NewScanner(scanPath, rulesPath)
		
		// Run scan
		results, err := sc.Scan()
		if err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}
		
		// Filter by severity if specified
		if severityFilter != "" {
			results = results.FilterBySeverity(severityFilter)
		}
		
		// Apply compliance framework filtering
		if complianceFramework != "" {
			results = results.FilterByCompliance(complianceFramework)
		}
		
		// Generate report
		reporter := reports.NewReporter(outputFormat, outputPath)
		if err := reporter.Generate(results); err != nil {
			return fmt.Errorf("report generation failed: %w", err)
		}
		
		// Print summary
		summary := results.Summary()
		fmt.Printf("\n✅ Scan Complete!\n")
		fmt.Printf("   Critical: %d\n", summary.Critical)
		fmt.Printf("   High: %d\n", summary.High)
		fmt.Printf("   Medium: %d\n", summary.Medium)
		fmt.Printf("   Low: %d\n", summary.Low)
		fmt.Printf("   Warnings: %d\n", summary.Warnings)
		fmt.Printf("\n📊 Report saved to: %s\n", outputPath)
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	
	scanCmd.Flags().StringVarP(&complianceFramework, "compliance", "c", "", "Compliance framework (soc2, pci-dss, hipaa)")
	scanCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format (json, pdf)")
	scanCmd.Flags().StringVarP(&outputPath, "output", "o", "security-report", "Output file path (without extension)")
	scanCmd.Flags().StringVarP(&rulesPath, "rules", "r", "rules", "Path to rules directory")
	scanCmd.Flags().StringVarP(&severityFilter, "severity", "s", "", "Filter by severity (critical, high, medium, low)")
}

func Execute() error {
	return rootCmd.Execute()
}

