package imap

import (
	"context"
	"sync"
)

type IMAPService struct {
	clients      map[string]*client.Client
	configs      map[string]interfaces.MailboxConfig
	eventHandler func(interfaces.MailEvent)
	clientsMutex sync.RWMutex
	wg           sync.WaitGroup
	ctx          context.Context
	cancel       context.CancelFunc
	statuses     map[string]interfaces.MailboxStatus
	statusMutex  sync.RWMutex
}
