package reports

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/iac-security-scanner/scanner/types"
)

type Reporter struct {
	format string
	outputPath string
}

func NewReporter(format, outputPath string) *Reporter {
	return &Reporter{
		format: format,
		outputPath: outputPath,
	}
}

func (r *Reporter) Generate(results *types.ScanResults) error {
	switch r.format {
	case "json":
		return r.generateJSON(results)
	case "pdf":
		return r.generatePDF(results)
	default:
		return fmt.Errorf("unsupported format: %s", r.format)
	}
}

func (r *Reporter) generateJSON(results *types.ScanResults) error {
	outputFile := r.outputPath
	if !endsWith(outputFile, ".json") {
		outputFile += ".json"
	}
	
	report := Report{
		Timestamp:   time.Now().Format(time.RFC3339),
		Summary:     results.Summary(),
		Issues:      results.Issues,
		FilesScanned: results.Files,
	}
	
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	
	return nil
}

func (r *Reporter) generatePDF(results *types.ScanResults) error {
	outputFile := r.outputPath
	if !endsWith(outputFile, ".pdf") {
		outputFile += ".pdf"
	}
	
	pdfGen := NewPDFGenerator()
	return pdfGen.Generate(results, outputFile)
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

type Report struct {
	Timestamp    string           `json:"timestamp"`
	Summary      types.Summary    `json:"summary"`
	Issues       []types.Issue    `json:"issues"`
	FilesScanned []string         `json:"files_scanned"`
}

