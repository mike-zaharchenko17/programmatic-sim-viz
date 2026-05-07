package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"psv-generator/internal/generator"
	"sync"
	"syscall"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func RunServer() {
	// this context applies to everything (server AND pipeline)
	// a cancel here shuts the whole thing down
	ctx, cancel := context.WithCancel(context.Background())

	// gracefully handle interrupts (Render sends SIGTERM on shutdown)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	auctionResultChan := make(chan generator.AuctionResult)
	pipelineWindDownChan := make(chan struct{}, 1)

	pipelineMgr := PipelineManager{
		mu:        sync.Mutex{},
		isRunning: false,
		serverCtx: ctx,
		outCh:     auctionResultChan,
		stopCh:    pipelineWindDownChan,
	}

	bh := BroadcastHub{
		PipelineMgr: &pipelineMgr,
		ChannelMap:  make(map[chan generator.AuctionResult]bool),
		mu:          sync.RWMutex{},
	}

	go bh.Run()

	// lazy start
	if len(bh.ChannelMap) > 0 {
		pipelineMgr.StartIfNeeded()
	}

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/ws", WsHandlerWithHub(&bh))

	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	sc := echo.StartConfig{Address: ":" + port}

	if err := sc.Start(ctx, e); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
