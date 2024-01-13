package entity

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"

// Deprecated: use neo4jutil.NodeLabel* instead
const (
	NodeLabelContact             = neo4jutil.NodeLabelContact
	NodeLabelJobRole             = neo4jutil.NodeLabelJobRole
	NodeLabelEmail               = neo4jutil.NodeLabelEmail
	NodeLabelLocation            = neo4jutil.NodeLabelLocation
	NodeLabelInteractionEvent    = neo4jutil.NodeLabelInteractionEvent
	NodeLabelInteractionSession  = neo4jutil.NodeLabelInteractionSession
	NodeLabelNote                = neo4jutil.NodeLabelNote
	NodeLabelLogEntry            = neo4jutil.NodeLabelLogEntry
	NodeLabelOrganization        = neo4jutil.NodeLabelOrganization
	NodeLabelBillingProfile      = neo4jutil.NodeLabelBillingProfile
	NodeLabelMasterPlan          = neo4jutil.NodeLabelMasterPlan
	NodeLabelMasterPlanMilestone = neo4jutil.NodeLabelMasterPlanMilestone
	NodeLabelAction              = neo4jutil.NodeLabelAction
	NodeLabelPageView            = neo4jutil.NodeLabelPageView
	NodeLabelPhoneNumber         = neo4jutil.NodeLabelPhoneNumber
	NodeLabelTag                 = neo4jutil.NodeLabelTag
	NodeLabelIssue               = neo4jutil.NodeLabelIssue
	NodeLabelUser                = neo4jutil.NodeLabelUser
	NodeLabelAnalysis            = neo4jutil.NodeLabelAnalysis
	NodeLabelAttachment          = neo4jutil.NodeLabelAttachment
	NodeLabelMeeting             = neo4jutil.NodeLabelMeeting
	NodeLabelSocial              = neo4jutil.NodeLabelSocial
	NodeLabelPlayer              = neo4jutil.NodeLabelPlayer
	NodeLabelCountry             = neo4jutil.NodeLabelCountry
	NodeLabelComment             = neo4jutil.NodeLabelComment
	NodeLabelServiceLineItem     = neo4jutil.NodeLabelServiceLineItem
	NodeLabelOpportunity         = neo4jutil.NodeLabelOpportunity
	NodeLabelInvoicingCycle      = neo4jutil.NodeLabelInvoicingCycle
)
