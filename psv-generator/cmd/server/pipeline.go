package server

import (
	"context"
	"psv-generator/internal/generator"
	"sync"
)

func RunPipeline(
	ctx context.Context,
	auctionResultChan chan generator.AuctionResult,
	pipelineWindDownChan chan struct{},
) {
	pipelineCtx, cancel := context.WithCancel(ctx)

	go func() {
		select {
		case <-ctx.Done():
			return
		case <-pipelineWindDownChan:
			cancel()
		}
	}()

	bidRequestChan := make(chan *generator.BidRequest)
	bidResponseChan := make(chan []*generator.BidResponse)

	wg := sync.WaitGroup{}

	wg.Add(3)

	go generator.BidRequestProducer(pipelineCtx, bidRequestChan, &wg)
	go generator.BidRequestResponsePipe(pipelineCtx, bidRequestChan, bidResponseChan, &wg)
	go generator.BidResponseConsumer(pipelineCtx, bidResponseChan, auctionResultChan, &wg)

	wg.Wait()
}
