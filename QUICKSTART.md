# Quick Start Guide

Get up and running with the Infrastructure-as-Code Security Scanner in 5 minutes.

## Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile commands)

## Installation

### Option 1: Build from Source

```bash
cd iac-security-scanner
go mod download
make build
# or: go build -o iac-audit .
```

### Option 2: Use Makefile

```bash
make deps    # Install dependencies
make build   # Build the binary
```

## First Scan

### 1. Test with Example Files

```bash
# Scan the included examples
./iac-audit scan ./examples

# Or with Makefile
make run-example
```

### 2. Scan Your Own Infrastructure

```bash
# Scan Terraform files
./iac-audit scan ./terraform

# Scan with HIPAA compliance check
./iac-audit scan ./terraform --compliance=hipaa

# Generate PDF report
./iac-audit scan ./terraform --format=pdf --output=my-security-report
```

## Understanding the Output

### Console Output

```
🔍 Scanning: ./terraform
📋 Compliance: hipaa
📄 Output: security-report (json)

✅ Scan Complete!
   Critical: 3
   High: 5
   Medium: 8
   Low: 2
   Warnings: 1

📊 Report saved to: security-report.json
```

### JSON Report Structure

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "summary": {
    "critical": 3,
    "high": 5,
    "medium": 8,
    "low": 2,
    "warnings": 1,
    "total": 19
  },
  "issues": [
    {
      "id": "aws-s3-public-access-1",
      "rule_id": "aws-s3-public-access",
      "severity": "critical",
      "title": "S3 Bucket Public Access",
      "description": "S3 bucket allows public access...",
      "file": "./terraform/main.tf",
      "compliance": ["soc2", "pci-dss", "hipaa"],
      "remediation": {
        "description": "Remove public access...",
        "steps": ["Step 1", "Step 2"],
        "code_example": "resource \"aws_s3_bucket\"..."
      }
    }
  ],
  "files_scanned": ["./terraform/main.tf"]
}
```

## Common Use Cases

### 1. Pre-Commit Security Check

```bash
# Scan only critical and high severity issues
./iac-audit scan . --severity=critical
./iac-audit scan . --severity=high
```

### 2. Compliance Audit

```bash
# SOC 2 compliance check
./iac-audit scan ./infrastructure --compliance=soc2 --format=pdf

# PCI-DSS compliance check
./iac-audit scan ./infrastructure --compliance=pci-dss --format=pdf
```

### 3. CI/CD Integration

```bash
# In your CI pipeline
./iac-audit scan ./terraform --format=json --output=ci-report.json

# Check exit code (non-zero if critical issues found)
if [ $? -ne 0 ]; then
  echo "Security issues found!"
  exit 1
fi
```

## Customizing Rules

### Add Your Own Rules

1. Create a YAML file in the `rules/` directory:

```yaml
id: my-custom-rule
name: My Custom Security Rule
description: Checks for specific security issue
severity: high
type: terraform
compliance:
  - soc2
conditions:
  - resource: aws_s3_bucket
    property: versioning
    operator: missing
remediation:
  description: Enable versioning
  steps:
    - Add versioning configuration
  code_example: |
    resource "aws_s3_bucket" "example" {
      versioning {
        enabled = true
      }
    }
```

2. Use custom rules directory:

```bash
./iac-audit scan ./terraform --rules=./my-custom-rules
```

## Troubleshooting

### Issue: "No rules found"

**Solution**: Ensure the rules directory exists and contains YAML files:
```bash
ls -la rules/
```

### Issue: "Failed to parse file"

**Solution**: The file might not be valid Terraform/CloudFormation/K8s. Check the file format.

### Issue: "No issues found"

**Solution**: This is good! But if you expect issues, check:
- Are your rules correctly defined?
- Does your infrastructure match the rule conditions?
- Try scanning the example files first

## Next Steps

- Read the [full README](README.md) for detailed documentation
- Explore the [example rules](rules/) to understand rule structure
- Check out [example infrastructure](examples/) for testing
- Customize rules for your organization's needs

## Getting Help

- Check the [README.md](README.md) for detailed documentation
- Review example rules in the `rules/` directory
- Test with the example files in `examples/`

## Building for Distribution

```bash
# Build for all platforms
make build-all

# Or manually
GOOS=linux GOARCH=amd64 go build -o iac-audit-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o iac-audit-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o iac-audit-windows-amd64.exe .
```

The binaries will be standalone and can be distributed without Go installed.

