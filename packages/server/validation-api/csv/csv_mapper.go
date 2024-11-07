package csv

import (
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/validation-api/model"
)

func GenerateCSVRow(data model.ValidateEmailMailSherpaData) (header []string, row []string) {
	// Define the header fields (without the excluded ones)
	header = []string{
		"Email", "SyntaxIsValid", "User", "Domain", "CleanEmail",
		"IsFirewalled", "Provider", "SecureGatewayProvider", "IsCatchAll", "CanConnectSMTP",
		"HasMXRecord", "HasSPFRecord", "TLSRequired", "IsPrimaryDomain", "PrimaryDomain",
		"SkippedValidation", "Deliverable", "IsMailboxFull", "IsRoleAccount", "IsSystemGenerated",
		"IsFreeAccount", "SmtpSuccess", "RetryValidation", "TLSRequired", "AlternateEmail",
	}

	// Create the corresponding row for the given data
	row = []string{
		data.Email,
		fmt.Sprintf("%v", data.Syntax.IsValid),
		data.Syntax.User,
		data.Syntax.Domain,
		data.Syntax.CleanEmail,
		fmt.Sprintf("%v", data.DomainData.IsFirewalled),
		data.DomainData.Provider,
		data.DomainData.SecureGatewayProvider,
		fmt.Sprintf("%v", data.DomainData.IsCatchAll),
		fmt.Sprintf("%v", data.DomainData.CanConnectSMTP),
		fmt.Sprintf("%v", data.DomainData.HasMXRecord),
		fmt.Sprintf("%v", data.DomainData.HasSPFRecord),
		fmt.Sprintf("%v", data.DomainData.TLSRequired),
		fmt.Sprintf("%v", data.DomainData.IsPrimaryDomain),
		data.DomainData.PrimaryDomain,
		fmt.Sprintf("%v", data.EmailData.SkippedValidation),
		data.EmailData.Deliverable,
		fmt.Sprintf("%v", data.EmailData.IsMailboxFull),
		fmt.Sprintf("%v", data.EmailData.IsRoleAccount),
		fmt.Sprintf("%v", data.EmailData.IsSystemGenerated),
		fmt.Sprintf("%v", data.EmailData.IsFreeAccount),
		fmt.Sprintf("%v", data.EmailData.SmtpSuccess),
		fmt.Sprintf("%v", data.EmailData.RetryValidation),
		fmt.Sprintf("%v", data.EmailData.TLSRequired),
		data.EmailData.AlternateEmail,
	}

	return header, row
}
