package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

func BidRequestProducer(ctx context.Context, bidRequestChan chan *BidRequest, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			close(bidRequestChan)
			return
		case bidRequestChan <- generateBidRequest():
			time.Sleep(time.Second)
		}
	}
}

// consumer of BidRequest, producer of BidResponse- intermediate step
func BidRequestResponsePipe(ctx context.Context, bidRequestChan chan *BidRequest, bidResponseChan chan []*BidResponse, wg *sync.WaitGroup) {
	var x *BidRequest
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case x = <-bidRequestChan:
			data, _ := json.MarshalIndent(&x, "", "   ")
			fmt.Println("Received data on bidRequestChannel:\n" + string(data))
			bidResponseChan <- generateBidResponses(x)
		}
	}
}

func BidResponseConsumer(ctx context.Context, bidResponseChan chan []*BidResponse, wg *sync.WaitGroup) {
	var resArray []*BidResponse

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case resArray = <-bidResponseChan:
			data, _ := json.MarshalIndent(&resArray, "", "  ")
			fmt.Println("Received data on bidResponseChannel:\n" + string(data))
		}
	}
}
