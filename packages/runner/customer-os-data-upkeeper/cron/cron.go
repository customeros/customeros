package cron

import (
	"sync"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
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
	GroupGlobalContact   = "globalContact"
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
	GroupEmail                             = "email"
	GroupEmailBulk                         = "emailBulkValidation"
	GroupMailstack                         = "mailstack"
	GroupSendEmails                        = "sendEmails"
	GroupProcessEmails                     = "processSentEmails"
	GroupRampMailboxes                     = "rampUpMailboxes"
	GroupIngestEmailsFromProvidersRealtime = "ingestEmailsFromProvidersRealtime"
	GroupIngestEmailsFromProvidersHistory  = "ingestEmailsFromProvidersHistory"
	GroupIngestEmailsSendToAgents          = "ingestEmailsSendToAgents"

	// Flow related groups
	GroupFlow      = "flowExecutionGroup"
	GroupFlowStats = "flowStatisticsGroup"

	// Other groups
	GroupDomain         = "domain"
	GroupReminder       = "reminder"
	GroupWebSession     = "webSession"
	GroupTouchpoint     = "refreshLastTouchpoint"
	GroupUnthreadIssues = "linkUnthreadIssues"
	GroupTenant         = "tenant"
	GroupAgent          = "agent"
	GroupScraper        = "scraper"
)

// LOCK MANAGEMENT
var jobLocks = struct {
	sync.Mutex
	locks map[string]*sync.Mutex
}{
	locks: map[string]*sync.Mutex{
		GroupOrganization:                      {},
		GroupGlobalOrg:                         {},
		GroupContact:                           {},
		GroupGlobalContact:                     {},
		GroupContactBetter:                     {},
		GroupLinkedInAsk:                       {},
		GroupLinkedInProcess:                   {},
		GroupContactEnrich:                     {},
		GroupOrphanContacts:                    {},
		GroupContract:                          {},
		GroupInvoice:                           {},
		GroupTouchpoint:                        {},
		GroupCurrency:                          {},
		GroupUnthreadIssues:                    {},
		GroupEmail:                             {},
		GroupEmailBulk:                         {},
		GroupFlow:                              {},
		GroupFlowStats:                         {},
		GroupRampMailboxes:                     {},
		GroupIngestEmailsFromProvidersRealtime: {},
		GroupIngestEmailsFromProvidersHistory:  {},
		GroupIngestEmailsSendToAgents:          {},
		GroupSendEmails:                        {},
		GroupProcessEmails:                     {},
		GroupDomain:                            {},
		GroupMailstack:                         {},
		GroupReminder:                          {},
		GroupWebSession:                        {},
		GroupTenant:                            {},
		GroupAgent:                             {},
		GroupScraper:                           {},
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
	addJob(cont.Cfg.App.Cron.CronScheduleIcpCheck, GroupOrganization, findLeads, "findLeads")

	addJob(cont.Cfg.App.Cron.CronScheduleSyncDataToGlobalOrgs, GroupGlobalOrg, syncDataToGlobalOrgs, "syncDataToGlobalOrgs")
	addJob(cont.Cfg.App.Cron.CronScheduleProcessWebsiteForGlobalOrgs, GroupGlobalOrg, processWebsiteForGlobalOrgs, "processWebsiteForGlobalOrgs")
	addJob(cont.Cfg.App.Cron.CronScheduleEnrichGlobalOrg, GroupGlobalOrg, enrichGlobalOrganization, "enrichGlobalOrganization")
	addJob(cont.Cfg.App.Cron.CronScheduleSyncFromGlobalOrgsToTenantOrgs, GroupGlobalOrg, syncGlobalOrgsToTenantOrganizations, "syncGlobalOrgsToTenantOrganizations")
	addJob(cont.Cfg.App.Cron.CronScheduleGlobalOrgScrape, GroupScraper, scrapeGlobalOrganizations, "scrapeGlobalOrganizations")
	addJob(cont.Cfg.App.Cron.CronScheduleLinkExtractionFromScrapedPage, GroupScraper, extractPageLinks, "extractPageLinks")
	addJob(cont.Cfg.App.Cron.CronScheduleDownloadIconAndLogo, GroupGlobalOrg, downloadGlobalOrganizationLogo, "downloadGlobalOrganizationLogo")

	// Contract Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleUpdateContract, GroupContract, updateContractsStatusAndRenewal, "updateContractsStatusAndRenewal")

	// Invoice Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleGenerateInvoice, GroupInvoice, generateCycleInvoices, "generateCycleInvoices")
	addJob(cont.Cfg.App.Cron.CronScheduleGenerateOffCycleInvoice, GroupInvoice, generateOffCycleInvoices, "generateOffCycleInvoices")
	addJob(cont.Cfg.App.Cron.CronScheduleGenerateNextPreviewInvoice, GroupInvoice, generateNextPreviewInvoices, "generateNextPreviewInvoices")
	addJob(cont.Cfg.App.Cron.CronScheduleCleanupInvoices, GroupInvoice, cleanupInvoices, "cleanupInvoices")
	addJob(cont.Cfg.App.Cron.CronScheduleUpkeepInvoices, GroupInvoice, upkeepInvoices, "upkeepInvoices")
	addJob(cont.Cfg.App.Cron.CronScheduleAdjustInvoiceStatus, GroupInvoice, adjustInvoiceStatus, "adjustInvoiceStatus")
	addJob(cont.Cfg.App.Cron.CronScheduleSendPayInvoiceNotification, GroupInvoice, sendPayInvoiceNotifications, "sendPayInvoiceNotifications")
	addJob(cont.Cfg.App.Cron.CronScheduleSendRemindInvoiceNotification, GroupInvoice, sendRemindInvoiceNotifications, "sendRemindInvoiceNotifications")

	// Contact Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleUpkeepContacts, GroupContact, upkeepContacts, "upkeepContacts")
	addJob(cont.Cfg.App.Cron.CronScheduleEnrichWithWorkEmailFromBetterContact, GroupContactBetter, enrichWithWorkEmailFromBetterContactJob, "enrichWithWorkEmailFromBetterContact")
	addJob(cont.Cfg.App.Cron.CronScheduleCheckBetterContactRequestsWithoutResponse, GroupContactBetter, checkBetterContactRequestsWithoutResponseJob, "checkBetterContactRequestsWithoutResponse")
	addJob(cont.Cfg.App.Cron.CronScheduleAskForLinkedInConnections, GroupLinkedInAsk, askForLinkedInConnections, "askForLinkedInConnections")
	addJob(cont.Cfg.App.Cron.CronScheduleProcessLinkedInConnections, GroupLinkedInProcess, processLinkedInConnections, "processLinkedInConnections")
	addJob(cont.Cfg.App.Cron.CronScheduleEnrichContacts, GroupContactEnrich, enrichContacts, "enrichContacts")
	addJob(cont.Cfg.App.Cron.CronScheduleLinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn, GroupOrphanContacts, linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn, "linkOrphanContacts")

	// Global Contact Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleSyncDataToGlobalContacts, GroupGlobalContact, syncDataToGlobalContacts, "syncDataToGlobalContacts")
	addJob(cont.Cfg.App.Cron.CronScheduleDownloadContactProfilePhoto, GroupGlobalContact, downloadContactProfilePhoto, "downloadContactProfilePhoto")
	addJob(cont.Cfg.App.Cron.CronScheduleEnrichGlobalContactWithBettercontact, GroupGlobalContact, enrichGlobalOrganizationWithBettercontact, "enrichGlobalOrganizationWithBettercontact")

	// Email Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleValidateEmails, GroupEmail, validateEmails, "validateEmails")
	addJob(cont.Cfg.App.Cron.CronScheduleValidateEmailsFromBulkRequests, GroupEmailBulk, validateEmailsFromBulkRequests, "validateEmailsFromBulkRequests")
	addJob(cont.Cfg.App.Cron.CronScheduleCheckScrubbyResult, GroupEmail, checkScrubbyResult, "checkScrubbyResult")
	addJob(cont.Cfg.App.Cron.CronScheduleCheckEnrowResults, GroupEmail, checkEnrowResult, "checkEnrowResult")
	addJob(cont.Cfg.App.Cron.CronScheduleCleanEmails, GroupEmail, cleanEmails, "cleanEmails")
	addJob(cont.Cfg.App.Cron.CronScheduleSendEmails, GroupSendEmails, sendEmails, "sendEmails")
	addJob(cont.Cfg.App.Cron.CronScheduleProcessSentEmails, GroupProcessEmails, processSentEmails, "processSentEmails")

	addJob(cont.Cfg.App.Cron.CronScheduleIngestEmailsFromProviders, GroupIngestEmailsFromProvidersRealtime, ingestEmailsFromProvidersRealtime, "ingestEmailsFromProvidersRealtime")
	addJob(cont.Cfg.App.Cron.CronScheduleIngestEmailsFromProviders, GroupIngestEmailsFromProvidersHistory, ingestEmailsFromProvidersHistory, "ingestEmailsFromProvidersHistory")
	addJob(cont.Cfg.App.Cron.CronScheduleIngestEmailsSendToAgents, GroupIngestEmailsSendToAgents, ingestEmailsSendToAgents, "ingestEmailsSendToAgents")

	// Flow Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleFlowExecution, GroupFlow, flowExecution, "flowExecution")
	addJob(cont.Cfg.App.Cron.CronScheduleFlowStatistics, GroupFlowStats, flowStatistics, "flowStatistics")
	addJob(cont.Cfg.App.Cron.CronScheduleRampUpMailboxes, GroupRampMailboxes, rampUpMailboxes, "rampUpMailboxes")

	// Tenant Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleCheckTenantOnboarding, GroupTenant, checkTenantOnboarding, "checkTenantOnboarding")

	// Agent Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleRerunAgent, GroupAgent, rerunAgent, "rerunAgent")

	// Other Jobs
	addJob(cont.Cfg.App.Cron.CronScheduleRefreshLastTouchpoint, GroupTouchpoint, refreshLastTouchpoint, "refreshLastTouchpoint")
	addJob(cont.Cfg.App.Cron.CronScheduleGetCurrencyRatesECB, GroupCurrency, getCurrencyRatesECB, "getCurrencyRatesECB")
	addJob(cont.Cfg.App.Cron.CronScheduleLinkUnthreadIssues, GroupUnthreadIssues, linkUnthreadIssues, "linkUnthreadIssues")
	addJob(cont.Cfg.App.Cron.CronScheduleCheckDomains, GroupDomain, checkDomains, "checkDomains")
	addJob(cont.Cfg.App.Cron.CronScheduleMailstackReputation, GroupMailstack, checkMailstackDomainReputation, "checkMailstackDomainReputation")
	addJob(cont.Cfg.App.Cron.CronScheduleSendOrganizationsReminders, GroupReminder, sendReminders, "sendReminders")
	addJob(cont.Cfg.App.Cron.CronScheduleProcessWebSessions, GroupWebSession, processWebSessions, "processWebSessions")
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
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).UpkeepOrganizations()
}

func findLeads(cont *container.Container) {
	cont.AgentProducers.NewLeadProducer.Execute()
}

func syncDataToGlobalOrgs(cont *container.Container) {
	service.NewGlobalOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).SyncDataIntoGlobalOrganizations()
}

func processWebsiteForGlobalOrgs(cont *container.Container) {
	service.NewGlobalOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).ScrapinCompanyByWebsite()
}

// Contract Jobs
func updateContractsStatusAndRenewal(cont *container.Container) {
	service.NewContractService(cont.Cfg, cont.Log, cont.Repositories, cont.CommonServices).UpkeepContracts()
}

// Invoice Jobs
func generateCycleInvoices(cont *container.Container) {
	if cont.Cfg.App.ProcessConfig.CycleInvoicingEnabled == false {
		cont.Log.Warn("INVOICING IS DISABLED")
		return
	}
	cont.AgentProducers.InvoiceProducer.Execute()
}

func generateOffCycleInvoices(cont *container.Container) {
	// service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories).GenerateOffCycleInvoices()
}

func generateNextPreviewInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories).GenerateNextPreviewInvoices()
}

func cleanupInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories).CleanupInvoices()
}

func upkeepInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories).UpkeepInvoices()
}

func adjustInvoiceStatus(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.CommonServices, cont.Repositories).AdjustInvoiceStatus()
}

func sendPayInvoiceNotifications(cont *container.Container) {
	cont.AgentProducers.SendInvoiceProducer.Execute()
}

func sendRemindInvoiceNotifications(cont *container.Container) {
	cont.AgentProducers.SendInvoiceProducer.Execute()
}

// Contact Jobs
func upkeepContacts(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices).UpkeepContacts()
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

func ingestEmailsFromProvidersRealtime(cont *container.Container) {
	service.NewIngestEmailService(cont.Cfg, cont.Log, cont.CommonServices).SyncEmailsInState(postgres_entity.REAL_TIME)
}

func ingestEmailsFromProvidersHistory(cont *container.Container) {
	service.NewIngestEmailService(cont.Cfg, cont.Log, cont.CommonServices).SyncEmailsInState(postgres_entity.HISTORY)
}

func ingestEmailsSendToAgents(cont *container.Container) {
	service.NewIngestEmailService(cont.Cfg, cont.Log, cont.CommonServices).SendIngestedEmailsToAgents()
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
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).RefreshLastTouchpoint()
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
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).SendReminders()
}

func processWebSessions(cont *container.Container) {
	cont.AgentProducers.NewWebSessionProducer.Execute()
}

func syncGlobalOrgsToTenantOrganizations(cont *container.Container) {
	service.NewGlobalOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).SyncGlobalOrgsToTenantOrganizations()
}

func scrapeGlobalOrganizations(cont *container.Container) {
	service.NewGlobalOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).ScrapeGlobalOrgs()
}

func downloadGlobalOrganizationLogo(cont *container.Container) {
	service.NewMediaService(cont.Log, cont.CommonServices).FetchAndStoreCompanyLogos()
}

func downloadContactProfilePhoto(cont *container.Container) {
	service.NewMediaService(cont.Log, cont.CommonServices).FetchAndStoreContactProfilePhotos()
}

func checkTenantOnboarding(cont *container.Container) {
	service.NewTenantService(cont.Cfg, cont.Log, cont.CommonServices).CheckOnboarding()
}

func rerunAgent(cont *container.Container) {
	service.NewAgentService(cont.Cfg, cont.Log, cont.CommonServices).RerunExecutions()
}

func extractPageLinks(cont *container.Container) {
	service.NewGlobalOrganizationService(cont.Cfg, cont.Log, cont.CommonServices).ExtractWebpageLinks()
}

func syncDataToGlobalContacts(cont *container.Container) {
	service.NewGlobalContactService(cont.Cfg, cont.Log, cont.CommonServices).SyncDataIntoGlobalContacts()
}

func enrichGlobalOrganizationWithBettercontact(cont *container.Container) {
	service.NewGlobalContactService(cont.Cfg, cont.Log, cont.CommonServices).EnrichGlobalOrganizationWithBettercontact()
}
