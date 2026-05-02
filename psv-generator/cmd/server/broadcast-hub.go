package server

import (
	"fmt"
	"psv-generator/internal/generator"
	"sync"
	"time"
)

const idleTimeout = 45 * time.Second

type BroadcastHub struct {
	PipelineMgr *PipelineManager
	// array of subscribed channels
	ChannelMap map[chan generator.AuctionResult]bool

	idleTimer *time.Timer

	mu sync.RWMutex
}

func (hub *BroadcastHub) Subscribe(clientChannel chan generator.AuctionResult) {
	hub.mu.Lock()
	hub.ChannelMap[clientChannel] = true
	hub.cancelIdle()
	hub.mu.Unlock()

	hub.PipelineMgr.StartIfNeeded()
}

func (hub *BroadcastHub) Unsubscribe(clientChannel chan generator.AuctionResult) {
	hub.mu.Lock()
	delete(hub.ChannelMap, clientChannel)
	close(clientChannel)
	if len(hub.ChannelMap) == 0 {
		hub.scheduleIdle()
	}
	hub.mu.Unlock()
}

// whenever there is data generated from the pipeline, broadcast it to all
// subscribed clients. If send is not successful, log it

func (hub *BroadcastHub) Run() {
	for dataForBroadcast := range hub.PipelineMgr.outCh {
		hub.mu.RLock() // "I am broadcasting nobody touch this data"
		for clientCh, subscribed := range hub.ChannelMap {
			if subscribed {
				select {
				case clientCh <- dataForBroadcast:
				default:
					// placeholder
					fmt.Printf("did not send to channel")
				}
			}
		}
		hub.mu.RUnlock()
	}
}

func (hub *BroadcastHub) scheduleIdle() {
	if hub.idleTimer != nil {
		if !hub.idleTimer.Stop() {
			select {
			case <-hub.idleTimer.C:
			default:
			}
		}
	}

	hub.idleTimer = time.AfterFunc(idleTimeout, func() {
		hub.mu.RLock()
		empty := len(hub.ChannelMap) == 0
		hub.mu.RUnlock()
		if empty {
			hub.PipelineMgr.stopCh <- struct{}{}
		}
	})
}

func (hub *BroadcastHub) cancelIdle() {
	if hub.idleTimer == nil {
		return
	}

	if !hub.idleTimer.Stop() {
		select {
		case <-hub.idleTimer.C:
		default:
		}
	}

	hub.idleTimer = nil
}
