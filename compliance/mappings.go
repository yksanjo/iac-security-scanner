package compliance

import (
	"github.com/iac-security-scanner/scanner/types"
)

// MapRuleToCompliance maps security rules to compliance frameworks
var ComplianceMappings = map[string][]types.ComplianceFramework{
	// SOC 2 mappings
	"encryption-at-rest":     {types.ComplianceSOC2},
	"encryption-in-transit":  {types.ComplianceSOC2},
	"access-control":         {types.ComplianceSOC2},
	"audit-logging":          {types.ComplianceSOC2},
	"data-backup":            {types.ComplianceSOC2},
	
	// PCI-DSS mappings
	"encryption-at-rest":     {types.CompliancePCIDSS},
	"encryption-in-transit":  {types.CompliancePCIDSS},
	"no-public-access":      {types.CompliancePCIDSS},
	"secure-storage":         {types.CompliancePCIDSS},
	"network-segmentation":   {types.CompliancePCIDSS},
	
	// HIPAA mappings
	"encryption-at-rest":     {types.ComplianceHIPAA},
	"encryption-in-transit":  {types.ComplianceHIPAA},
	"access-control":         {types.ComplianceHIPAA},
	"audit-logging":          {types.ComplianceHIPAA},
	"data-retention":         {types.ComplianceHIPAA},
}

func GetComplianceForRule(ruleID string) []types.ComplianceFramework {
	if frameworks, exists := ComplianceMappings[ruleID]; exists {
		return frameworks
	}
	return []types.ComplianceFramework{}
}

