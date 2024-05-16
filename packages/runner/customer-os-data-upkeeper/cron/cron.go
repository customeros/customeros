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
	invoiceGroup               = "invoice"
	refreshLastTouchpointGroup = "refreshLastTouchpoint"
	currencyGroup              = "currency"
	linkUnthreadIssuesGroup    = "linkUnthreadIssues"
)

var jobLocks = struct {
	sync.Mutex
	locks map[string]*sync.Mutex
}{
	locks: map[string]*sync.Mutex{
		organizationGroup:          {},
		contractGroup:              {},
		invoiceGroup:               {},
		refreshLastTouchpointGroup: {},
		currencyGroup:              {},
		linkUnthreadIssuesGroup:    {},
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
		cont.Log.Fatalf("Could not add cron job %s: %v", "updateContractsStatusAndRenewal", err.Error())
	}

	err = c.AddFunc(cont.Cfg.Cron.CronScheduleWebScrapeOrganization, func() {
		lockAndRunJob(cont, organizationGroup, webScrapeOrganizations)
	})
	if err != nil {
		cont.Log.Fatalf("Could not add cron job %s: %v", "webScrapeOrganizations", err.Error())
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
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).UpkeepOrganizations()
}

func webScrapeOrganizations(cont *container.Container) {
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).WebScrapeOrganizations()
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
	service.NewOrganizationService(cont.Cfg, cont.Log, cont.Repositories, cont.EventProcessingServicesClient).RefreshLastTouchpoint()
}

func getCurrencyRatesECB(cont *container.Container) {
	service.NewCurrencyService(cont.Cfg, cont.Log, cont.Repositories).GetCurrencyRatesECB()
}

func linkUnthreadIssues(cont *container.Container) {
	service.NewIssueService(cont.Cfg, cont.Log, cont.Repositories).LinkUnthreadIssues()
}
