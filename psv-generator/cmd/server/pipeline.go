package server

import (
	"context"
	"psv-generator/internal/generator"
	"sync"
)

type PipelineManager struct {
	mu        sync.Mutex
	isRunning bool
	cancel    context.CancelFunc

	serverCtx context.Context
	outCh     chan generator.AuctionResult
	stopCh    chan struct{}
}

func (m *PipelineManager) StartIfNeeded() {
	m.mu.Lock()

	// prevent run if the pipeline is already running
	if m.isRunning {
		m.mu.Unlock()
		return
	}

	drainStop(m.stopCh)

	runCtx, runCancel := context.WithCancel(m.serverCtx)
	m.cancel = runCancel
	m.isRunning = true
	m.mu.Unlock()

	go func() {
		select {
		case <-m.serverCtx.Done():
			m.cancel()
		case <-m.stopCh:
			m.cancel()
		}
	}()

	go func() {
		defer func() {
			m.mu.Lock()
			m.isRunning = false
			m.cancel = nil
			m.mu.Unlock()
		}()

		runPipeline(runCtx, m.outCh)
	}()

}

// drains a channel without blocking
func drainStop(ch chan struct{}) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}

func runPipeline(
	runCtx context.Context,
	auctionResultChan chan generator.AuctionResult,
) {
	bidRequestChan := make(chan *generator.BidRequest)
	bidResponseChan := make(chan []*generator.BidResponse)

	wg := sync.WaitGroup{}

	wg.Add(3)

	go generator.BidRequestProducer(runCtx, bidRequestChan, &wg)
	go generator.BidRequestResponsePipe(runCtx, bidRequestChan, bidResponseChan, &wg)
	go generator.BidResponseConsumer(runCtx, bidResponseChan, auctionResultChan, &wg)

	wg.Wait()
}
