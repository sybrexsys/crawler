package manager

import (
	"context"
	"sync"

	"github.com/sybrexsys/crawler/config"
	"github.com/sybrexsys/crawler/logs"
)

type Manager struct {
	log             *logs.Logger
	cfg             *config.Config
	sitesMutex      sync.RWMutex
	sites           map[string]*siteManager
	totalProcessing uint32
	C               chan struct{}
	E               chan error
	wg              *sync.WaitGroup
}

func NewManager(ctx context.Context, log *logs.Logger, cfg *config.Config, startUrl string, wg *sync.WaitGroup) *Manager {
	tmp := &Manager{
		log:   log,
		cfg:   cfg,
		sites: make(map[string]*siteManager),
		C:     make(chan struct{}),
		E:     make(chan error),
		wg:    wg,
	}
	tmp.AddUrl(ctx, startUrl)
	return tmp
}

func (mgr *Manager) AddUrl(ctx context.Context, url string) {
	t := newSiteManager(mgr, 10)
	go t.Run(ctx, mgr.wg)
}
