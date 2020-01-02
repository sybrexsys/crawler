package manager

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

func loadData(ctx context.Context, url string) error {
	return nil
}

func newWorker(ctx context.Context, mngr *Manager, workerID int, list chan<- int, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		list <- workerID
	}()
	mngr.log.Printf("Worker [%d] started", workerID)
	rand.Seed(time.Now().UnixNano())
	d := time.Second * time.Duration(rand.Int31n(20))
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
	mngr.log.Printf("Worker [%d] finished", workerID)
}
