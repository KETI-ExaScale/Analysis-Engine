package main

import (
	"analysis-engine/pkg/analysis"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	w := analysis.InitEngine()
	wg.Add(2)
	w.Work(ctx, &wg)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	cancel()
	wg.Wait()
	os.Exit(0)
}
