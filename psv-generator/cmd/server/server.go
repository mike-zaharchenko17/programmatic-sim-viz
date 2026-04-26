package server

import (
	"context"
	"psv-generator/internal/generator"
	"sync"
	"time"
)

func RunPipeline() {
	ctx, cancel := context.WithCancel(context.Background())

	bidRequestChan := make(chan *generator.BidRequest)
	bidResponseChan := make(chan []*generator.BidResponse)
	auctionResultChan := make(chan generator.AuctionResult)

	wg := sync.WaitGroup{}

	wg.Add(2)

	go generator.BidRequestProducer(ctx, bidRequestChan, &wg)
	go generator.BidRequestResponsePipe(ctx, bidRequestChan, bidResponseChan, &wg)
	go generator.BidResponseConsumer(ctx, bidResponseChan, auctionResultChan, &wg)

	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()

	wg.Wait()
}

func RunServer() {
	RunPipeline()
}
