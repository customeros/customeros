package cron

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/container"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/logger"
	"github.com/openline-ai/openline-customer-os/packages/runner/customer-os-data-upkeeper/service"
	"github.com/robfig/cron"
	"sync"
)

const (
	organizationGroup = "organization"
	contractGroup     = "contract"
)

var jobLocks = struct {
	sync.Mutex
	locks map[string]*sync.Mutex
}{
	locks: map[string]*sync.Mutex{
		organizationGroup: {},
		contractGroup:     {},
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
