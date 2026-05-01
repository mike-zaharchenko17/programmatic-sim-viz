package server

import (
	"context"
	"os"
	"psv-generator/internal/generator"
	"sync"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RunPipeline(
	ctx context.Context,
	auctionResultChan chan generator.AuctionResult,
	pipelineWindDownChan chan struct{},
) {
	pipelineCtx, cancel := context.WithCancel(ctx)

	go func() {
		<-pipelineWindDownChan
		cancel()
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

func RunServer() {
	// this context applies to everything (server AND pipeline)
	// a cancel here shuts the whole thing down
	ctx, cancel := context.WithCancel(context.Background())

	// gracefully handle interrupts
	sigCh := make(chan os.Signal, 1)
	go func() {
		<-sigCh
		cancel()
	}()

	auctionResultChan := make(chan generator.AuctionResult)
	pipelineWindDownChan := make(chan struct{}, 1)

	bh := BroadcastHub{
		SourceChannel:           auctionResultChan,
		ChannelMap:              make(map[chan generator.AuctionResult]bool),
		PipelineWindDownChannel: pipelineWindDownChan,
		mu:                      sync.RWMutex{},
	}

	go bh.Run()

	go RunPipeline(ctx, auctionResultChan, pipelineWindDownChan)

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Static("/", "../public")

	e.GET("/ws", WsHandlerWithHub(&bh))

	sc := echo.StartConfig{Address: ":1323"}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
