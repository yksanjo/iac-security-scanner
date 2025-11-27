package reports

import (
	"fmt"
	"strings"

	"github.com/iac-security-scanner/scanner/types"
	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator struct{}

func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{}
}

func (p *PDFGenerator) Generate(results *types.ScanResults, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Infrastructure-as-Code Security Report")
	pdf.Ln(12)
	
	// Summary
	summary := results.Summary()
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Summary")
	pdf.Ln(8)
	
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 8, fmt.Sprintf("Total Issues: %d", summary.Total))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("Critical: %d", summary.Critical))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("High: %d", summary.High))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("Medium: %d", summary.Medium))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("Low: %d", summary.Low))
	pdf.Ln(6)
	pdf.Cell(40, 8, fmt.Sprintf("Warnings: %d", summary.Warnings))
	pdf.Ln(10)
	
	// Files Scanned
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Files Scanned")
	pdf.Ln(8)
	
	pdf.SetFont("Arial", "", 10)
	for _, file := range results.Files {
		pdf.Cell(40, 6, file)
		pdf.Ln(5)
	}
	pdf.Ln(5)
	
	// Issues
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Security Issues")
	pdf.Ln(8)
	
	pdf.SetFont("Arial", "", 10)
	for i, issue := range results.Issues {
		if i > 0 {
			pdf.Ln(5)
		}
		
		// Issue header
		pdf.SetFont("Arial", "B", 11)
		pdf.SetTextColor(255, 0, 0) // Red for severity
		pdf.Cell(40, 8, fmt.Sprintf("[%s] %s", strings.ToUpper(string(issue.Severity)), issue.Title))
		pdf.Ln(6)
		
		pdf.SetTextColor(0, 0, 0) // Black
		pdf.SetFont("Arial", "", 10)
		
		// Description
		pdf.Cell(40, 6, fmt.Sprintf("File: %s", issue.File))
		pdf.Ln(5)
		pdf.MultiCell(170, 6, issue.Description, "", "", false)
		pdf.Ln(3)
		
		// Compliance
		if len(issue.Compliance) > 0 {
			compStr := "Compliance: " + strings.Join(func() []string {
				var comps []string
				for _, c := range issue.Compliance {
					comps = append(comps, string(c))
				}
				return comps
			}(), ", ")
			pdf.Cell(40, 6, compStr)
			pdf.Ln(5)
		}
		
		// Remediation
		if issue.Remediation.Description != "" {
			pdf.SetFont("Arial", "B", 10)
			pdf.Cell(40, 6, "Remediation:")
			pdf.Ln(5)
			pdf.SetFont("Arial", "", 9)
			pdf.MultiCell(170, 5, issue.Remediation.Description, "", "", false)
			
			if len(issue.Remediation.Steps) > 0 {
				pdf.Ln(2)
				for j, step := range issue.Remediation.Steps {
					pdf.Cell(10, 5, "")
					pdf.Cell(40, 5, fmt.Sprintf("%d. %s", j+1, step))
					pdf.Ln(4)
				}
			}
			
			if issue.Remediation.CodeExample != "" {
				pdf.Ln(2)
				pdf.SetFont("Courier", "", 8)
				pdf.MultiCell(170, 4, issue.Remediation.CodeExample, "", "", false)
				pdf.SetFont("Arial", "", 10)
			}
		}
		
		// Add page break if needed
		if pdf.GetY() > 250 {
			pdf.AddPage()
		}
	}
	
	return pdf.OutputFileAndClose(outputPath)
}

