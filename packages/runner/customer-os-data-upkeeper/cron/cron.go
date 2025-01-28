package cron

import (
	"sync"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/robfig/cron"

	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/container"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/customeros/customeros/packages/runner/customer-os-data-upkeeper/service"
)

// CONSTANTS - Group definitions
const (
	// Organization related groups
	GroupOrganization = "organization"
	GroupGlobalOrg    = "globalOrganization"

	// Contact related groups
	GroupContact         = "contact"
	GroupContactEnrich   = "contactEnrich"
	GroupContactBetter   = "contactEnrichWithBettercontact"
	GroupLinkedInAsk     = "askForLinkedInConnections"
	GroupLinkedInProcess = "processLinkedInConnections"
	GroupOrphanContacts  = "orphanContactsGroup"

	// Financial related groups
	GroupContract = "contract"
	GroupInvoice  = "invoice"
	GroupCurrency = "currency"

	// Email related groups
	GroupEmail         = "email"
	GroupEmailBulk     = "emailBulkValidation"
	GroupMailstack     = "mailstack"
	GroupSendEmails    = "sendEmails"
	GroupProcessEmails = "processSentEmails"
	GroupRampMailboxes = "rampUpMailboxes"

	// Flow related groups
	GroupFlow      = "flowExecutionGroup"
	GroupFlowStats = "flowStatisticsGroup"

	// Other groups
	GroupDomain         = "domain"
	GroupReminder       = "reminder"
	GroupWebSession     = "webSession"
	GroupTouchpoint     = "refreshLastTouchpoint"
	GroupUnthreadIssues = "linkUnthreadIssues"
)

// LOCK MANAGEMENT
var jobLocks = struct {
	sync.Mutex
	locks map[string]*sync.Mutex
}{
	locks: map[string]*sync.Mutex{
		GroupOrganization:    {},
		GroupGlobalOrg:       {},
		GroupContact:         {},
		GroupContactBetter:   {},
		GroupLinkedInAsk:     {},
		GroupLinkedInProcess: {},
		GroupContactEnrich:   {},
		GroupOrphanContacts:  {},
		GroupContract:        {},
		GroupInvoice:         {},
		GroupTouchpoint:      {},
		GroupCurrency:        {},
		GroupUnthreadIssues:  {},
		GroupEmail:           {},
		GroupEmailBulk:       {},
		GroupFlow:            {},
		GroupFlowStats:       {},
		GroupRampMailboxes:   {},
		GroupSendEmails:      {},
		GroupProcessEmails:   {},
		GroupDomain:          {},
		GroupMailstack:       {},
		GroupReminder:        {},
		GroupWebSession:      {},
	},
}

// CORE FUNCTIONALITY
func StartCron(cont *container.Container) *cron.Cron {
	c := cron.New()
	registerJobs(c, cont)
	c.Start()
	return c
}

func StopCron(log logger.Logger, cron *cron.Cron) error {
	log.Info("Gracefully stopping cron")
	cron.Stop()
	return nil
}

// JOB REGISTRATION
func registerJobs(c *cron.Cron, cont *container.Container) {
	// Helper function to reduce duplication
	addJob := func(schedule string, group string, job func(*container.Container), name string) {
		err := c.AddFunc(schedule, func() {
			lockAndRunJob(cont, group, job)
		})
		if err != nil {
			cont.Log.Fatalf("Could not add cron job %s: %v", name, err.Error())
		}
	}

	// Organization Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleUpdateOrganization, GroupOrganization, updateOrganizations, "updateOrganizations")
	addJob(cont.Cfg.App.Cron.CronScheduleSyncDataToGlobalOrgs, GroupGlobalOrg, syncDataToGlobalOrgs, "syncDataToGlobalOrgs")
	addJob(cont.Cfg.App.Cron.CronScheduleProcessWebsiteForGlobalOrgs, GroupGlobalOrg, processWebsiteForGlobalOrgs, "processWebsiteForGlobalOrgs")
	addJob(cont.Cfg.App.Cron.CronScheduleEnrichGlobalOrg, GroupGlobalOrg, enrichGlobalOrganization, "enrichGlobalOrganization")
	addJob(cont.Cfg.App.Cron.CronScheduleSyncFromGlobalOrgsToTenantOrgs, GroupGlobalOrg, syncGlobalOrgsToTenantOrganizations, "syncGlobalOrgsToTenantOrganizations")

	// Contract Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleUpdateContract, GroupContract, updateContractsStatusAndRenewal, "updateContractsStatusAndRenewal")

	// Invoice Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleGenerateInvoice, GroupInvoice, generateCycleInvoices, "generateCycleInvoices")
	addJob(cont.Cfg.App.Cron.CronScheduleGenerateOffCycleInvoice, GroupInvoice, generateOffCycleInvoices, "generateOffCycleInvoices")
	addJob(cont.Cfg.App.Cron.CronScheduleGenerateNextPreviewInvoice, GroupInvoice, generateNextPreviewInvoices, "generateNextPreviewInvoices")
	addJob(cont.Cfg.App.Cron.CronScheduleGenerateInvoicePaymentLink, GroupInvoice, generateInvoicePaymentLinks, "generateInvoicePaymentLinks")
	addJob(cont.Cfg.App.Cron.CronScheduleCheckInvoiceFinalized, GroupInvoice, sendInvoiceFinalizedEvents, "sendInvoiceFinalizedEvents")
	addJob(cont.Cfg.App.Cron.CronScheduleCleanupInvoices, GroupInvoice, cleanupInvoices, "cleanupInvoices")
	addJob(cont.Cfg.App.Cron.CronScheduleAdjustInvoiceStatus, GroupInvoice, adjustInvoiceStatus, "adjustInvoiceStatus")
	addJob(cont.Cfg.App.Cron.CronScheduleSendPayInvoiceNotification, GroupInvoice, sendPayInvoiceNotifications, "sendPayInvoiceNotifications")
	addJob(cont.Cfg.App.Cron.CronScheduleSendRemindInvoiceNotification, GroupInvoice, sendRemindInvoiceNotifications, "sendRemindInvoiceNotifications")

	// Contact Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleUpkeepContacts, GroupContact, upkeepContacts, "upkeepContacts")
	addJob(cont.Cfg.App.Cron.CronScheduleAskForWorkEmailOnBetterContact, GroupContactBetter, askForWorkEmailOnBetterContactJob, "askForWorkEmailOnBetterContact")
	addJob(cont.Cfg.App.Cron.CronScheduleEnrichWithWorkEmailFromBetterContact, GroupContactBetter, enrichWithWorkEmailFromBetterContactJob, "enrichWithWorkEmailFromBetterContact")
	addJob(cont.Cfg.App.Cron.CronScheduleCheckBetterContactRequestsWithoutResponse, GroupContactBetter, checkBetterContactRequestsWithoutResponseJob, "checkBetterContactRequestsWithoutResponse")
	addJob(cont.Cfg.App.Cron.CronScheduleAskForLinkedInConnections, GroupLinkedInAsk, askForLinkedInConnections, "askForLinkedInConnections")
	addJob(cont.Cfg.App.Cron.CronScheduleProcessLinkedInConnections, GroupLinkedInProcess, processLinkedInConnections, "processLinkedInConnections")
	addJob(cont.Cfg.App.Cron.CronScheduleEnrichContacts, GroupContactEnrich, enrichContacts, "enrichContacts")
	addJob(cont.Cfg.App.Cron.CronScheduleLinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn, GroupOrphanContacts, linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn, "linkOrphanContacts")

	// Email Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleValidateEmails, GroupEmail, validateEmails, "validateEmails")
	addJob(cont.Cfg.App.Cron.CronScheduleValidateEmailsFromBulkRequests, GroupEmailBulk, validateEmailsFromBulkRequests, "validateEmailsFromBulkRequests")
	addJob(cont.Cfg.App.Cron.CronScheduleCheckScrubbyResult, GroupEmail, checkScrubbyResult, "checkScrubbyResult")
	addJob(cont.Cfg.App.Cron.CronScheduleCheckEnrowResults, GroupEmail, checkEnrowResult, "checkEnrowResult")
	addJob(cont.Cfg.App.Cron.CronScheduleCleanEmails, GroupEmail, cleanEmails, "cleanEmails")
	addJob(cont.Cfg.App.Cron.CronScheduleSendEmails, GroupSendEmails, sendEmails, "sendEmails")
	addJob(cont.Cfg.App.Cron.CronScheduleProcessSentEmails, GroupProcessEmails, processSentEmails, "processSentEmails")

	// Flow Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleFlowExecution, GroupFlow, flowExecution, "flowExecution")
	addJob(cont.Cfg.App.Cron.CronScheduleFlowStatistics, GroupFlowStats, flowStatistics, "flowStatistics")
	addJob(cont.Cfg.App.Cron.CronScheduleRampUpMailboxes, GroupRampMailboxes, rampUpMailboxes, "rampUpMailboxes")

	// Other Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleRefreshLastTouchpoint, GroupTouchpoint, refreshLastTouchpoint, "refreshLastTouchpoint")
	addJob(cont.Cfg.App.Cron.CronScheduleGetCurrencyRatesECB, GroupCurrency, getCurrencyRatesECB, "getCurrencyRatesECB")
	addJob(cont.Cfg.App.Cron.CronScheduleLinkUnthreadIssues, GroupUnthreadIssues, linkUnthreadIssues, "linkUnthreadIssues")
	addJob(cont.Cfg.App.Cron.CronScheduleCheckDomains, GroupDomain, checkDomains, "checkDomains")
	addJob(cont.Cfg.App.Cron.CronScheduleMailstackReputation, GroupMailstack, checkMailstackDomainReputation, "checkMailstackDomainReputation")
	addJob(cont.Cfg.App.Cron.CronScheduleSendOrganizationsReminders, GroupReminder, sendReminders, "sendReminders")
	addJob(cont.Cfg.App.Cron.CronScheduleProcessWebSessions, GroupWebSession, processWebSessions, "processWebSessions")
	addJob(cont.Cfg.App.Cron.CronScheduleAnalyzeWebSessionIntent, GroupWebSession, analyzeWebSessionIntent, "analyzeWebSessionIntent")
}

// HELPER FUNCTIONS
func lockAndRunJob(cont *container.Container, groupName string, job func(*container.Container)) {
	jobLocks.locks[groupName].Lock()
	defer jobLocks.locks[groupName].Unlock()
	defer tracing.RecoverAndLogToJaeger(cont.Log)
	job(cont)
}

// JOB IMPLEMENTATIONS
// Organization Jobs
func updateOrganizations(cont *container.Container) {
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.CommonServices, cont.EventProcessingServicesClient).UpkeepOrganizations()
}

func syncDataToGlobalOrgs(cont *container.Container) {
	service.NewGlobalOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).SyncDataIntoGlobalOrganizations()
}

func processWebsiteForGlobalOrgs(cont *container.Container) {
	service.NewGlobalOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).ScrapinCompanyByWebsite()
}

// Contract Jobs
func updateContractsStatusAndRenewal(cont *container.Container) {
	service.NewContractService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient, cont.CommonServices).UpkeepContracts()
}

// Invoice Jobs
func generateCycleInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories, cont.EventProcessingServicesClient).GenerateCycleInvoices()
}

func generateOffCycleInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories, cont.EventProcessingServicesClient).GenerateOffCycleInvoices()
}

func generateNextPreviewInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories, cont.EventProcessingServicesClient).GenerateNextPreviewInvoices()
}

func generateInvoicePaymentLinks(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories, cont.EventProcessingServicesClient).GenerateInvoicePaymentLinks()
}

func sendInvoiceFinalizedEvents(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories, cont.EventProcessingServicesClient).SendInvoiceFinalizedEvent()
}

func cleanupInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories, cont.EventProcessingServicesClient).CleanupInvoices()
}

func adjustInvoiceStatus(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories, cont.EventProcessingServicesClient).AdjustInvoiceStatus()
}

func sendPayInvoiceNotifications(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories, cont.EventProcessingServicesClient).SendPayNotifications()
}

func sendRemindInvoiceNotifications(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories, cont.EventProcessingServicesClient).SendRemindNotifications()
}

// Contact Jobs
func upkeepContacts(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices).UpkeepContacts()
}

func askForWorkEmailOnBetterContactJob(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices).AskForWorkEmailOnBetterContact()
}

func enrichWithWorkEmailFromBetterContactJob(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices).EnrichWithWorkEmailFromBetterContact()
}

func checkBetterContactRequestsWithoutResponseJob(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices).CheckBetterContactRequestsWithoutResponse()
}

func askForLinkedInConnections(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices).AskForLinkedInConnections()
}

func processLinkedInConnections(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices).ProcessLinkedInConnections()
}

func enrichContacts(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices).EnrichContacts()
}

func linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices).LinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn()
}

// Email Jobs
func validateEmails(cont *container.Container) {
	service.NewEmailService(cont.Cfg, cont.Log, cont.CommonServices).ValidateEmails()
}

func validateEmailsFromBulkRequests(cont *container.Container) {
	service.NewEmailService(cont.Cfg, cont.Log, cont.CommonServices).ValidateEmailsFromBulkRequests()
}

func checkScrubbyResult(cont *container.Container) {
	service.NewEmailService(cont.Cfg, cont.Log, cont.CommonServices).CheckScrubbyResult()
}

func checkEnrowResult(cont *container.Container) {
	service.NewEmailService(cont.Cfg, cont.Log, cont.CommonServices).CheckEnrowRequestsWithoutResponse()
}

func cleanEmails(cont *container.Container) {
	service.NewEmailService(cont.Cfg, cont.Log, cont.CommonServices).CleanEmails()
}

func sendEmails(cont *container.Container) {
	service.NewEmailService(cont.Cfg, cont.Log, cont.CommonServices).SendEmails()
}

func processSentEmails(cont *container.Container) {
	service.NewEmailService(cont.Cfg, cont.Log, cont.CommonServices).ProcessSentEmails()
}

// Flow Jobs
func flowExecution(cont *container.Container) {
	service.NewFlowExecutionService(cont.Cfg, cont.Log, cont.CommonServices).ExecuteScheduledFlowActions()
}

func flowStatistics(cont *container.Container) {
	service.NewFlowExecutionService(cont.Cfg, cont.Log, cont.CommonServices).ComputeFlowStatistics()
}

func rampUpMailboxes(cont *container.Container) {
	service.NewFlowExecutionService(cont.Cfg, cont.Log, cont.CommonServices).RampUpMailboxes()
}

// Other Jobs
func refreshLastTouchpoint(cont *container.Container) {
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.CommonServices, cont.EventProcessingServicesClient).RefreshLastTouchpoint()
}

func getCurrencyRatesECB(cont *container.Container) {
	service.NewCurrencyService(cont.Cfg, cont.Log, cont.Repositories).GetCurrencyRatesECB()
}

func linkUnthreadIssues(cont *container.Container) {
	service.NewIssueService(cont.Cfg, cont.Log, cont.Repositories).LinkUnthreadIssues()
}

func checkDomains(cont *container.Container) {
	service.NewDomainService(cont.Cfg, cont.Log, cont.CommonServices).CheckDomains()
}

func checkMailstackDomainReputation(cont *container.Container) {
	service.NewMailstackService(cont.Cfg, cont.Log, cont.CommonServices).CheckMailstackDomainReputation()
}

func enrichGlobalOrganization(cont *container.Container) {
	service.NewGlobalOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).EnrichGlobalOrganization()
}

func sendReminders(cont *container.Container) {
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.CommonServices, cont.EventProcessingServicesClient).SendReminders()
}

func processWebSessions(cont *container.Container) {
	service.NewWebSessionService(cont.Cfg, cont.Log, cont.CommonServices).ProcessWebSessions()
}

func analyzeWebSessionIntent(cont *container.Container) {
	service.NewWebSessionService(cont.Cfg, cont.Log, cont.CommonServices).ProcessIntentSignals()
}

func syncGlobalOrgsToTenantOrganizations(cont *container.Container) {
	service.NewGlobalOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).SyncGlobalOrgsToTenantOrganizations()
}
