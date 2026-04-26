package server

import (
	"fmt"
	"psv-generator/internal/generator"
	"sync"
)

type BroadcastHub struct {
	// master channel
	SourceChannel <-chan generator.AuctionResult
	// array of subscribed channels
	ChannelMap map[chan generator.AuctionResult]bool

	mu sync.RWMutex
}

func (hub *BroadcastHub) Subscribe(clientChannel chan generator.AuctionResult) {
	hub.mu.Lock()
	hub.ChannelMap[clientChannel] = true
	hub.mu.Unlock()
}

func (hub *BroadcastHub) Unsubscribe(clientChannel chan generator.AuctionResult) {
	hub.mu.Lock()
	delete(hub.ChannelMap, clientChannel)
	close(clientChannel)
	hub.mu.Unlock()
}

// whenever there is data generated from the pipeline, broadcast it to all
// subscribed clients. If send is not successful, log it

func (hub *BroadcastHub) Run() {
	for dataForBroadcast := range hub.SourceChannel {
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
