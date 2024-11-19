package service

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type EmailAnalyzer interface {
	Analyze(email EmailMessageData) HeaderAnalysis
}

type emailAnalyzer struct{}

func NewEmailAnalyzer() EmailAnalyzer {
	return &emailAnalyzer{}
}

func (a *emailAnalyzer) ProcessCheck(email EmailMessageData) HeaderAnalysis {
	analysis := HeaderAnalysis{
		ShouldProcess: true, // Default to processing
	}

	// Check auto-responder first
	if a.isAutoResponder(email.Headers) {
		analysis.IsAutoResponder = true
		analysis.ShouldProcess = false
		return analysis
	}

	// Check bounce
	if a.isBounce(email.Headers) {
		analysis.IsBounce = true
		analysis.ShouldProcess = false
		return analysis
	}

	// Check bulk mail
	if a.isBulkMail(email.Headers) {
		analysis.IsBulkMail = true
		analysis.ShouldProcess = false
		return analysis
	}

	return analysis
}

func (a *emailAnalyzer) isAutoResponder(headers EmailHeaders) bool {
	return headers.XAutoreply != "" ||
		headers.XAutoresponse != "" ||
		headers.AutoSubmitted ||
		headers.XLoop ||
		headers.Precedence == "auto_reply"
}

func (a *emailAnalyzer) isBounce(headers EmailHeaders) bool {
	return headers.XFailedRecepients ||
		headers.DeliveryStatus ||
		headers.ContentDescription == "delivery report" ||
		a.isReturnPathBounce(headers.ReturnPath)
}

func (a *emailAnalyzer) isBulkMail(headers EmailHeaders) bool {
	return headers.ListUnsubscribe ||
		headers.Precedence == "bulk"
}

func (a *emailAnalyzer) isReturnPathBounce(returnPath string) bool {
	return returnPath == "" ||
		returnPath == "mailer-daemon" ||
		returnPath == "postmaster"
}
