package cron

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/container"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/service"
	"github.com/robfig/cron"
	"sync"
)

const (
	organizationGroup          = "organization"
	contractGroup              = "contract"
	orphanContactsGroup        = "orphanContactsGroup"
	invoiceGroup               = "invoice"
	refreshLastTouchpointGroup = "refreshLastTouchpoint"
	currencyGroup              = "currency"
	linkUnthreadIssuesGroup    = "linkUnthreadIssues"
	contactGroup               = "contact"
	contactEnrichGroup         = "contactEnrich"
	contactBettercontactGroup  = "contactEnrichWithBettercontact"
	contactWeconnectGroup      = "contactWeconnect"
	apiCacheGroup              = "api_cache"
	workflowGroup              = "workflow"
	emailGroup                 = "email"
	emailBulkValidationGroup   = "emailBulkValidation"
)

var jobLocks = struct {
	sync.Mutex
	locks map[string]*sync.Mutex
}{
	locks: map[string]*sync.Mutex{
		organizationGroup:          {},
		contactGroup:               {},
		contactBettercontactGroup:  {},
		contactWeconnectGroup:      {},
		contactEnrichGroup:         {},
		orphanContactsGroup:        {},
		contractGroup:              {},
		invoiceGroup:               {},
		refreshLastTouchpointGroup: {},
		currencyGroup:              {},
		linkUnthreadIssuesGroup:    {},
		apiCacheGroup:              {},
		workflowGroup:              {},
		emailGroup:                 {},
		emailBulkValidationGroup:   {},
	},
}

func StartCron(cont *container.Container) *cron.Cron {
	c := cron.New()

	// Add jobs
	err := c.AddFunc(cont.Cfg.Cron.CronScheduleUpdateContract, func() {
		lockAndRunJob(cont, contractGroup, updateContractsStatusAndRenewal)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "updateContractsStatusAndRenewal", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleUpdateOrganization, func() {
		lockAndRunJob(cont, organizationGroup, updateOrganizations)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "updateOrganizations", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleGenerateInvoice, func() {
		lockAndRunJob(cont, invoiceGroup, generateCycleInvoices)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "generateCycleInvoices", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleGenerateOffCycleInvoice, func() {
		lockAndRunJob(cont, invoiceGroup, generateOffCycleInvoices)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "generateOffCycleInvoices", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleGenerateNextPreviewInvoice, func() {
		lockAndRunJob(cont, invoiceGroup, generateNextPreviewInvoices)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "generateNextPreviewInvoices", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleGenerateInvoicePaymentLink, func() {
		lockAndRunJob(cont, invoiceGroup, generateInvoicePaymentLinks)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "generateInvoicePaymentLinks", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleCheckInvoiceFinalized, func() {
		lockAndRunJob(cont, invoiceGroup, sendInvoiceFinalizedEvents)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "autoPayInvoices", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleCleanupInvoices, func() {
		lockAndRunJob(cont, invoiceGroup, cleanupInvoices)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "cleanupInvoices", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleAdjustInvoiceStatus, func() {
		lockAndRunJob(cont, invoiceGroup, adjustInvoiceStatus)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "adjustInvoiceStatus", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleSendPayInvoiceNotification, func() {
		lockAndRunJob(cont, invoiceGroup, sendPayInvoiceNotifications)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "sendPayInvoiceNotifications", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleRefreshLastTouchpoint, func() {
		lockAndRunJob(cont, refreshLastTouchpointGroup, refreshLastTouchpoint)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "refreshLastTouchpoint", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleGetCurrencyRatesECB, func() {
		lockAndRunJob(cont, currencyGroup, getCurrencyRatesECB)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "getCurrencyRatesECB", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleLinkUnthreadIssues, func() {
		lockAndRunJob(cont, linkUnthreadIssuesGroup, linkUnthreadIssues)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "linkUnthreadIssues", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleUpkeepContacts, func() {
		lockAndRunJob(cont, contactGroup, upkeepContacts)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "upkeepContacts", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleAskForWorkEmailOnBetterContact, func() {
		lockAndRunJob(cont, contactBettercontactGroup, askForWorkEmailOnBetterContactJob)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "askForWorkEmailOnBetterContactJob", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleEnrichWithWorkEmailFromBetterContact, func() {
		lockAndRunJob(cont, contactBettercontactGroup, enrichWithWorkEmailFromBetterContactJob)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "enrichWithWorkEmailFromBetterContactJob", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleCheckBetterContactRequestsWithoutResponse, func() {
		lockAndRunJob(cont, contactBettercontactGroup, checkBetterContactRequestsWithoutResponseJob)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "checkBetterContactRequestsWithoutResponseJob", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleWeConnectSyncContacts, func() {
		lockAndRunJob(cont, contactWeconnectGroup, weConnectContacts)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "weConnectContacts", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleEnrichContacts, func() {
		lockAndRunJob(cont, contactEnrichGroup, enrichContacts)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "enrichContacts", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleLinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn, func() {
		lockAndRunJob(cont, orphanContactsGroup, linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleRefreshApiCache, func() {
		lockAndRunJob(cont, apiCacheGroup, refreshApiCache)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "refreshApiCache", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleExecuteWorkflow, func() {
		lockAndRunJob(cont, workflowGroup, executeWorkflows)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "executeWorkflows", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleValidateCustomerOSEmails, func() {
		lockAndRunJob(cont, emailGroup, validateEmails)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "validateEmails", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleValidateEmailsFromBulkRequests, func() {
		lockAndRunJob(cont, emailBulkValidationGroup, validateEmailsFromBulkRequests)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "validateEmailFromBulkRequests", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleCheckScrubbyResult, func() {
		lockAndRunJob(cont, emailGroup, checkScrubbyResult)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "validateEmails", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleCheckEnrowResults, func() {
		lockAndRunJob(cont, emailGroup, checkEnrowResult)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "checkEnrowResult", err.Error())
	}

	c.Start()

	return c
}

func lockAndRunJob(cont *container.Container, groupName string, job func(cont *container.Container)) {
	jobLocks.locks[groupName].Lock()
	defer jobLocks.locks[groupName].Unlock()

	job(cont)
}

func StopCron(log logger.Logger, cron *cron.Cron) error {
	// Gracefully stop
	log.Info("Gracefully stopping cron")
	cron.Stop()
	return nil
}

func updateContractsStatusAndRenewal(cont *container.Container) {
	service.NewContractService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).UpkeepContracts()
}

func updateOrganizations(cont *container.Container) {
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.CommonServices, cont.EventProcessingServicesClient).UpkeepOrganizations()
}

func upkeepContacts(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices, cont.CustomerOSApiClient, cont.EventBufferStoreService).UpkeepContacts()
}

func askForWorkEmailOnBetterContactJob(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices, cont.CustomerOSApiClient, cont.EventBufferStoreService).AskForWorkEmailOnBetterContact()
}

func enrichWithWorkEmailFromBetterContactJob(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices, cont.CustomerOSApiClient, cont.EventBufferStoreService).EnrichWithWorkEmailFromBetterContact()
}

func checkBetterContactRequestsWithoutResponseJob(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices, cont.CustomerOSApiClient, cont.EventBufferStoreService).CheckBetterContactRequestsWithoutResponse()
}

func weConnectContacts(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices, cont.CustomerOSApiClient, cont.EventBufferStoreService).SyncWeConnectContacts()
}

func enrichContacts(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices, cont.CustomerOSApiClient, cont.EventBufferStoreService).EnrichContacts()
}

func linkOrphanContactsToOrganizationBaseOnLinkedinScrapIn(cont *container.Container) {
	service.NewContactService(cont.Cfg, cont.Log, cont.CommonServices, cont.CustomerOSApiClient, cont.EventBufferStoreService).LinkOrphanContactsToOrganizationBaseOnLinkedinScrapIn()
}

func generateCycleInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).GenerateCycleInvoices()
}

func generateOffCycleInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).GenerateOffCycleInvoices()
}

func generateNextPreviewInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).GenerateNextPreviewInvoices()
}

func generateInvoicePaymentLinks(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).GenerateInvoicePaymentLinks()
}

func sendInvoiceFinalizedEvents(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).SendInvoiceFinalizedEvent()
}

func cleanupInvoices(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).CleanupInvoices()
}

func adjustInvoiceStatus(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).AdjustInvoiceStatus()
}

func sendPayInvoiceNotifications(cont *container.Container) {
	service.NewInvoiceService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).SendPayNotifications()
}

func refreshLastTouchpoint(cont *container.Container) {
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.CommonServices, cont.EventProcessingServicesClient).RefreshLastTouchpoint()
}

func getCurrencyRatesECB(cont *container.Container) {
	service.NewCurrencyService(cont.Cfg, cont.Log, cont.Repositories).GetCurrencyRatesECB()
}

func linkUnthreadIssues(cont *container.Container) {
	service.NewIssueService(cont.Cfg, cont.Log, cont.Repositories).LinkUnthreadIssues()
}

func refreshApiCache(cont *container.Container) {
	service.NewApiCacheService(cont.Cfg, cont.Log, cont.Repositories, cont.CommonServices).RefreshApiCache()
}

func executeWorkflows(cont *container.Container) {
	service.NewWorkflowService(cont.Cfg, cont.Log, cont.Repositories, cont.CommonServices).ExecuteWorkflows()
}

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
