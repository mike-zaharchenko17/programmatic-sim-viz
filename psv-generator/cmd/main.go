package main

import (
	"context"
	"psv-generator/internal/generator"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	bidRequestChan := make(chan *generator.BidRequest)
	wg := sync.WaitGroup{}

	wg.Add(2)
	go generator.BidRequestProducer(ctx, bidRequestChan, &wg)
	go generator.BidRequestConsumer(ctx, bidRequestChan, &wg)

	go func() {
		time.Sleep(time.Second * 10)
		cancel()
	}()

	wg.Wait()
}
