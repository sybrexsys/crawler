package manager

import (
	"context"
	"sync"
	"time"
)

type siteManager struct {
	owner     *Manager
	workers   chan int
	newJob    chan struct{}
	pMutex    sync.RWMutex
	processed map[string]struct{}
	uMutex    sync.Mutex
	urls      []string
}

func newSiteManager(owner *Manager, size int) *siteManager {
	tmp := &siteManager{
		owner:   owner,
		workers: make(chan int, size),
		newJob:  make(chan struct{}, 1),
	}
	for i := 0; i < size; i++ {
		tmp.workers <- i
	}
	return tmp
}

func (sm *siteManager) notifyNewJob() {
	select {
	case sm.newJob <- struct{}{}:
	default:
	}
}

var cnt int

func (sm *siteManager) getURL() *string {
	s := "Test"
	sm.owner.log.Printf("New job with ID:%d was get", cnt)
	cnt++
	return &s
	/*
		sm.uMutex.Lock()
		defer sm.uMutex.Unlock()
		l := len(sm.urls)
		if l == 0 {
			return nil
		}
		str := sm.urls[0]
		copy(sm.urls[:l-1], sm.urls[1:l])
		sm.urls = sm.urls[:l-1]
		return &str*/
}

func (sm *siteManager) AddURL(url string) {
	sm.pMutex.RLock()
	_, ok := sm.processed[url]
	if ok {
		sm.pMutex.RUnlock()
		return
	}
	sm.pMutex.Lock()
	_, ok = sm.processed[url]
	if ok {
		sm.pMutex.Unlock()
		return
	}
	defer sm.pMutex.Unlock()
	sm.uMutex.Lock()
	defer sm.uMutex.Unlock()
	sm.urls = append(sm.urls, url)
	sm.processed[url] = struct{}{}
	sm.notifyNewJob()
}

func (sm *siteManager) Run(ctx context.Context, wg *sync.WaitGroup) {
	var wrk chan int
	var tick <-chan time.Time
	var l *limiter
	if sm.owner.cfg.MaximumRPM > 0 {
		l = newLimiter(sm.owner.cfg.MaximumRPM, time.Minute)
	}
	wrk = sm.workers
	for {
		select {
		case <-ctx.Done():
			return
		case workerID := <-wrk:
			if l != nil {
				d := l.check()
				if d != 0 {
					if d > time.Second*5 {
						sm.owner.log.Print("limit was arrived. Waiting ", d-d%time.Second)
					}
					tm := time.NewTimer(d)
					tick = tm.C
					wrk = nil
					sm.workers <- workerID
					continue
				}
			}
			url := sm.getURL()
			if url == nil {
				wrk = nil
				continue
			}
			wg.Add(1)
			go newWorker(ctx, sm.owner, workerID, sm.workers, wg)
			if l != nil {
				l.addNewTime()
			}
		case <-sm.newJob:
			wrk = sm.workers
		case <-tick:
			wrk = sm.workers
		}
	}
}
