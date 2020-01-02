package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sybrexsys/crawler/config"
	"github.com/sybrexsys/crawler/logs"
	"github.com/sybrexsys/crawler/manager"
)

func main() {
	flag.Parse()
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Printf("Cannot read the config file: %s", err.Error())
		os.Exit(255)
	}

	log, err := logs.NewLog(cfg)
	if err != nil {
		log.Printf("Cannot create a new logger: %s", err.Error())
		os.Exit(255)
	}
	defer func() {
		log.Info("logging stopped")
		log.Close()
	}()
	ctx, cancel := context.WithCancel(context.Background())

	stop := make(chan os.Signal, 1)

	var wg sync.WaitGroup
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	url := getFirstUrl()

	m := manager.NewManager(ctx, log, cfg, url, &wg)
	if err != nil {
		log.Errorf("invalid URL was detected:%s", url)
		os.Exit(255)
	}
	select {
	case <-m.C:
		log.Printf("Job finished successfully")
	case <-m.E:
		log.Errorf("Job finished with error: %s", err.Error())
	case cc := <-stop:
		cancel()
		log.Printf("Signal (%v) was detected. Process is being stopped", cc)
	}
	wg.Wait()

}

func getFirstUrl() string {
	return "test.com"
}
