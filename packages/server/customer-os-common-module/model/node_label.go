package model

import "strings"

const (
	NodeLabelTenant                    = "Tenant"
	NodeLabelTenantSettings            = "TenantSettings"
	NodeLabelTenantBillingProfile      = "TenantBillingProfile"
	NodeLabelBankAccount               = "BankAccount"
	NodeLabelTimelineEvent             = "TimelineEvent"
	NodeLabelContact                   = "Contact"
	NodeLabelCustomField               = "CustomField"
	NodeLabelJobRole                   = "JobRole"
	NodeLabelEmail                     = "Email"
	NodeLabelLocation                  = "Location"
	NodeLabelInteractionEvent          = "InteractionEvent"
	NodeLabelInteractionSession        = "InteractionSession"
	NodeLabelNote                      = "Note"
	NodeLabelLogEntry                  = "LogEntry"
	NodeLabelOrganization              = "Organization"
	NodeLabelBillingProfile            = "BillingProfile"
	NodeLabelMasterPlan                = "MasterPlan"
	NodeLabelMasterPlanMilestone       = "MasterPlanMilestone"
	NodeLabelAction                    = "Action"
	NodeLabelPageView                  = "PageView"
	NodeLabelPhoneNumber               = "PhoneNumber"
	NodeLabelTag                       = "Tag"
	NodeLabelIssue                     = "Issue"
	NodeLabelUser                      = "User"
	NodeLabelAttachment                = "Attachment"
	NodeLabelMeeting                   = "Meeting"
	NodeLabelSocial                    = "Social"
	NodeLabelPlayer                    = "Player"
	NodeLabelCountry                   = "Country"
	NodeLabelActionItem                = "ActionItem"
	NodeLabelComment                   = "Comment"
	NodeLabelContract                  = "Contract"
	NodeLabelDeletedContract           = "DeletedContract"
	NodeLabelDomain                    = "Domain"
	NodeLabelServiceLineItem           = "ServiceLineItem"
	NodeLabelOpportunity               = "Opportunity"
	NodeLabelInvoicingCycle            = "InvoicingCycle"
	NodeLabelExternalSystem            = "ExternalSystem"
	NodeLabelInvoice                   = "Invoice"
	NodeLabelInvoiceLine               = "InvoiceLine"
	NodeLabelOrganizationPlan          = "OrganizationPlan"
	NodeLabelOrganizationPlanMilestone = "OrganizationPlanMilestone"
	NodeLabelReminder                  = "Reminder"
	NodeLabelOffering                  = "Offering"
	NodeLabelFlow                      = "Flow"
	NodeLabelFlowAction                = "FlowAction"
	NodeLabelFlowParticipant           = "FlowParticipant"
	NodeLabelFlowActionSender          = "FlowActionSender"
	NodeLabelFlowExecutionSettings     = "FlowExecutionSettings"
	NodeLabelFlowActionExecution       = "FlowActionExecution"
)

func NodeLabelWithTenant(label string, tenant string) string {
	return label + "_" + tenant
}

func GetTenantFromLabels(labels []string, nodeLabel string) string {
	var result string

	for _, event := range labels {
		if strings.Index(event, nodeLabel+"_") == 0 {
			result = event
			break
		}
	}

	if result != "" {
		return result[len(nodeLabel)+1 : len(result)]
	} else {
		return ""
	}
}
