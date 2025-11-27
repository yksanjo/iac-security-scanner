# 🛡️ Infrastructure-as-Code Security Scanner

<div align="center">

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![Go](https://img.shields.io/badge/go-1.21+-00ADD8.svg)
![License](https://img.shields.io/badge/license-Proprietary-red.svg)
![Platform](https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey.svg)

**A comprehensive security auditing tool for Infrastructure-as-Code files**

*Scans Terraform, CloudFormation, and Kubernetes YAML for security misconfigurations and generates compliance reports*

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Documentation](#-documentation) • [Examples](#-examples)

</div>

---

## 📋 Table of Contents

- [Overview](#-overview)
- [Features](#-features)
- [Architecture](#-architecture)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Usage](#-usage)
- [Rule Definitions](#-rule-definitions)
- [Compliance Frameworks](#-compliance-frameworks)
- [Report Formats](#-report-formats)
- [Examples](#-examples)
- [Project Structure](#-project-structure)
- [Building](#-building-for-distribution)
- [Contributing](#-contributing)

---

## 🎯 Overview

The Infrastructure-as-Code Security Scanner is a powerful CLI tool that automatically scans your IaC files for security misconfigurations, compliance violations, and best practice violations. It supports multiple formats and generates detailed reports to help you secure your infrastructure before deployment.

```
┌─────────────────────────────────────────────────────────────┐
│                    IaC Security Scanner                      │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Terraform  ──┐                                              │
│               ├──►  Parser  ──►  Rule Engine  ──►  Report  │
│  CloudFormation ──┤                                          │
│               │                                              │
│  Kubernetes  ──┘                                              │
│                                                              │
│  Compliance: SOC 2 | PCI-DSS | HIPAA                         │
│  Output: JSON | PDF                                          │
└─────────────────────────────────────────────────────────────┘
```

### Why Use This Tool?

- 🔒 **Prevent Security Breaches** - Catch misconfigurations before deployment
- 📊 **Compliance Ready** - Built-in support for SOC 2, PCI-DSS, and HIPAA
- ⚡ **Fast & Offline** - No cloud dependencies, works completely offline
- 🎯 **Customizable** - Easy-to-write YAML rules for your specific needs
- 📦 **Standalone Binary** - Single executable, no runtime dependencies

---

## ✨ Features

### 🔍 Multi-Format Support
```
✅ Terraform (.tf, .tf.json)
✅ CloudFormation (YAML, JSON)
✅ Kubernetes (YAML manifests)
```

### 📋 Compliance Frameworks
```
✅ SOC 2 Type II
✅ PCI-DSS Level 1
✅ HIPAA
```

### 📊 Report Generation
```
✅ JSON Reports (machine-readable)
✅ PDF Reports (executive-ready)
✅ Detailed Remediation Steps
✅ Code Examples
```

### 🛡️ Security Checks

| Category | Examples |
|----------|----------|
| **Access Control** | Public S3 buckets, open security groups, public RDS |
| **Encryption** | Unencrypted EBS volumes, missing TLS |
| **Network Security** | Open ports, unrestricted CIDR blocks |
| **Resource Limits** | Missing K8s resource limits |
| **Best Practices** | Missing tags, improper IAM policies |

---

## 🏗️ Architecture

```
┌──────────────┐
│   CLI Tool   │  (Cobra Framework)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Scanner    │  (File Discovery & Orchestration)
└──────┬───────┘
       │
       ├──► Terraform Parser ──┐
       ├──► CloudFormation ────┼──► Rule Engine ──► Issue Generator
       └──► Kubernetes ────────┘
                              │
                              ▼
                       ┌──────────────┐
                       │   Reporter   │
                       └──────┬───────┘
                              │
                    ┌─────────┴─────────┐
                    │                   │
                    ▼                   ▼
              JSON Report          PDF Report
```

---

## 🚀 Installation

### Option 1: Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/iac-security-scanner.git
cd iac-security-scanner

# Install dependencies
go mod download

# Build the binary
go build -o iac-audit .

# Or use Makefile
make build
```

### Option 2: Download Pre-built Binary

Download the latest release for your platform:
- **Linux**: `iac-audit-linux-amd64`
- **macOS (Intel)**: `iac-audit-darwin-amd64`
- **macOS (Apple Silicon)**: `iac-audit-darwin-arm64`
- **Windows**: `iac-audit-windows-amd64.exe`

### Option 3: Using Makefile

```bash
make deps    # Install dependencies
make build   # Build binary
make build-all  # Build for all platforms
```

---

## ⚡ Quick Start

### 1. Scan Your Infrastructure

```bash
# Basic scan
./iac-audit scan ./terraform

# With compliance framework
./iac-audit scan ./terraform --compliance=hipaa

# Generate PDF report
./iac-audit scan ./terraform --format=pdf --output=security-report
```

### 2. View Results

```bash
🔍 Scanning: ./terraform
📋 Compliance: hipaa
📄 Output: security-report (pdf)

✅ Scan Complete!
   Critical: 3
   High: 5
   Medium: 8
   Low: 2
   Warnings: 1

📊 Report saved to: security-report.pdf
```

### 3. Test with Examples

```bash
# Scan included example files
./iac-audit scan ./examples --format=json
```

---

## 📖 Usage

### Command Syntax

```bash
iac-audit scan [path] [flags]
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--compliance` | `-c` | Compliance framework (soc2, pci-dss, hipaa) | - |
| `--format` | `-f` | Output format (json, pdf) | `json` |
| `--output` | `-o` | Output file path | `security-report` |
| `--rules` | `-r` | Path to rules directory | `rules` |
| `--severity` | `-s` | Filter by severity (critical, high, medium, low) | - |

### Common Use Cases

#### 1. Pre-Commit Security Check

```bash
# Check for critical issues only
./iac-audit scan . --severity=critical
```

#### 2. Compliance Audit

```bash
# SOC 2 compliance check with PDF report
./iac-audit scan ./infrastructure --compliance=soc2 --format=pdf
```

#### 3. CI/CD Integration

```bash
# In your CI pipeline
./iac-audit scan ./terraform --format=json --output=ci-report.json

# Fail build if critical issues found
if [ $? -ne 0 ]; then
  echo "❌ Security issues found!"
  exit 1
fi
```

#### 4. Custom Rules

```bash
# Use your own rule set
./iac-audit scan ./terraform --rules=./company-rules
```

---

## 📝 Rule Definitions

Rules are defined in YAML files in the `rules/` directory. Each rule specifies security checks and remediation steps.

### Rule Structure

```yaml
id: unique-rule-id
name: Human Readable Rule Name
description: Detailed description of what the rule checks
severity: critical|high|medium|low|warning
type: terraform|cloudformation|kubernetes
compliance:
  - soc2
  - pci-dss
  - hipaa
conditions:
  - resource: aws_s3_bucket
    property: acl
    operator: equals
    value: public-read
remediation:
  description: How to fix the issue
  steps:
    - Step 1: Description
    - Step 2: Description
  code_example: |
    resource "aws_s3_bucket" "example" {
      acl = "private"
    }
```

### Example Rule: S3 Public Access

```yaml
id: aws-s3-public-access
name: S3 Bucket Public Access
description: S3 bucket allows public access, violating security best practices
severity: critical
type: terraform
compliance:
  - soc2
  - pci-dss
  - hipaa
conditions:
  - resource: aws_s3_bucket
    property: acl
    operator: equals
    value: public-read
remediation:
  description: Remove public access and use IAM policies
  steps:
    - Set acl = "private"
    - Configure bucket policy for specific access
    - Enable S3 Block Public Access
  code_example: |
    resource "aws_s3_bucket" "secure" {
      bucket = "my-secure-bucket"
      acl    = "private"
    }
```

### Included Rules

| Rule ID | Description | Severity |
|---------|-------------|----------|
| `aws-s3-public-access` | S3 bucket with public access | Critical |
| `aws-ec2-no-encryption` | EC2 instance without encryption | High |
| `aws-rds-public-access` | Publicly accessible RDS instance | Critical |
| `k8s-no-resource-limits` | Kubernetes pod without resource limits | Medium |
| `cf-s3-public-read` | CloudFormation S3 without public access block | Critical |
| `cf-security-group-open-rdp` | Security group with open RDP access | Critical |

---

## 🏛️ Compliance Frameworks

### SOC 2 Type II

**Coverage:**
- ✅ Encryption at rest and in transit
- ✅ Access control and authentication
- ✅ Audit logging and monitoring
- ✅ Data backup and recovery
- ✅ Change management

**Usage:**
```bash
./iac-audit scan ./infrastructure --compliance=soc2
```

### PCI-DSS Level 1

**Coverage:**
- ✅ Encryption requirements (AES-256)
- ✅ Network segmentation
- ✅ Secure storage of cardholder data
- ✅ No public access to sensitive data
- ✅ Access control requirements

**Usage:**
```bash
./iac-audit scan ./infrastructure --compliance=pci-dss
```

### HIPAA

**Coverage:**
- ✅ PHI encryption requirements
- ✅ Access controls and authentication
- ✅ Audit logging (who, what, when)
- ✅ Data retention policies
- ✅ Breach notification requirements

**Usage:**
```bash
./iac-audit scan ./infrastructure --compliance=hipaa
```

---

## 📊 Report Formats

### JSON Report

Machine-readable format perfect for CI/CD integration and automation.

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
      "line": 5,
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

### PDF Report

Executive-ready PDF with:
- 📈 Executive summary with statistics
- 📋 Detailed issue descriptions
- 🔧 Step-by-step remediation guides
- 💻 Code examples for fixes
- 🏛️ Compliance framework mappings

**Sample PDF Structure:**
```
┌─────────────────────────────────────┐
│   Infrastructure Security Report   │
├─────────────────────────────────────┤
│                                     │
│  Summary                            │
│  ───────                            │
│  Critical: 3                        │
│  High: 5                            │
│  Medium: 8                           │
│  Total: 19                          │
│                                     │
│  Issues                             │
│  ──────                             │
│  [Issue 1]                          │
│  [Issue 2]                          │
│  ...                                │
│                                     │
│  Remediation Steps                  │
│  Code Examples                      │
└─────────────────────────────────────┘
```

---

## 💡 Examples

### Example 1: Full Infrastructure Scan

```bash
# Scan entire infrastructure directory
./iac-audit scan ./infrastructure \
  --compliance=soc2 \
  --format=pdf \
  --output=infrastructure-audit-2024
```

### Example 2: Kubernetes Security Check

```bash
# Scan Kubernetes manifests for high/critical issues
./iac-audit scan ./k8s/manifests \
  --severity=high \
  --format=json
```

### Example 3: Terraform Pre-Deployment Check

```bash
# Quick check before deployment
./iac-audit scan ./terraform \
  --severity=critical \
  --format=json \
  --output=pre-deploy-check
```

### Example 4: Custom Rules for Organization

```bash
# Use organization-specific rules
./iac-audit scan ./infrastructure \
  --rules=./company-security-rules \
  --compliance=hipaa \
  --format=pdf
```

---

## 📁 Project Structure

```
iac-security-scanner/
├── cmd/                    # CLI command definitions
│   └── root.go            # Main command and flags
├── scanner/                # Core scanning engine
│   └── scanner.go         # File discovery and orchestration
├── parsers/                # File format parsers
│   ├── terraform.go       # Terraform HCL parser
│   ├── cloudformation.go  # CloudFormation YAML/JSON parser
│   └── kubernetes.go      # Kubernetes YAML parser
├── rules/                  # Security rule definitions
│   ├── engine.go          # Rule evaluation engine
│   ├── aws-s3-public-access.yaml
│   ├── aws-ec2-no-encryption.yaml
│   ├── aws-rds-public-access.yaml
│   ├── k8s-no-resource-limits.yaml
│   ├── cf-s3-public-read.yaml
│   └── cf-security-group-open-rdp.yaml
├── reports/                # Report generators
│   ├── reporter.go        # Report interface
│   └── pdf.go             # PDF report generator
├── compliance/             # Compliance framework mappings
│   └── mappings.go        # SOC 2, PCI-DSS, HIPAA mappings
├── types/                  # Type definitions
│   └── types.go           # Core data structures
├── examples/               # Example IaC files for testing
│   ├── terraform/
│   ├── kubernetes/
│   └── cloudformation/
├── main.go                 # Entry point
├── go.mod                  # Go module dependencies
├── Makefile                # Build automation
├── README.md               # This file
├── QUICKSTART.md           # Quick start guide
└── LICENSE                 # License file
```

---

## 🔨 Building for Distribution

### Single Platform

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o iac-audit-linux-amd64 .

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o iac-audit-darwin-amd64 .

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o iac-audit-darwin-arm64 .

# Windows
GOOS=windows GOARCH=amd64 go build -o iac-audit-windows-amd64.exe .
```

### All Platforms (Makefile)

```bash
make build-all
```

This creates standalone binaries with no runtime dependencies - perfect for distribution!

---

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-rule
   ```
3. **Add your rules or improvements**
   - Add new security rules in `rules/`
   - Improve parsers or rule engine
   - Enhance documentation
4. **Submit a pull request**

### Adding New Rules

1. Create a YAML file in `rules/` directory
2. Follow the rule structure format
3. Test with example files
4. Submit PR with description

---

## 📄 License

**Proprietary** - All rights reserved

This software is proprietary and confidential. Unauthorized reproduction or distribution is prohibited.

For licensing inquiries: **licensing@iac-security-scanner.com**

---

## 💰 Pricing

| License Type | Price | Includes |
|-------------|-------|----------|
| **Individual** | $299 | Single developer license |
| **Team** | $999 | 5-seat license |
| **Enterprise** | $2,999 | Unlimited seats + source code |

**Purchase Options:**
- 🛒 [Gumroad](https://gumroad.com/iac-security-scanner)
- 🛒 [Lemon Squeezy](https://lemonsqueezy.com/iac-security-scanner)
- 📧 Direct sales: sales@iac-security-scanner.com

---

## 📞 Support

- 📧 **Email**: support@iac-security-scanner.com
- 📚 **Documentation**: [Full Docs](QUICKSTART.md)
- 🐛 **Issues**: [GitHub Issues](https://github.com/yourusername/iac-security-scanner/issues)

---

## 🗺️ Roadmap

- [ ] Additional cloud providers (Azure, GCP)
- [ ] More compliance frameworks (ISO 27001, NIST)
- [ ] CI/CD integration templates (GitHub Actions, GitLab CI)
- [ ] Custom rule builder UI
- [ ] Baseline comparison and trend analysis
- [ ] Real-time monitoring mode
- [ ] Integration with popular IDEs

---

## ⭐ Star History

If you find this project useful, please consider giving it a star! ⭐

---

<div align="center">

**Built with ❤️ for the DevOps and Security community**

[⬆ Back to Top](#-infrastructure-as-code-security-scanner)

</div>
