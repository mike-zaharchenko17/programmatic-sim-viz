package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
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
			delayFloat := rand.NormFloat64()*200 + 500
			delayFloat = max(50, min(2000, delayFloat))
			delay := int(math.Round(delayFloat))
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}
}

// consumer of BidRequest, producer of Bid Response- intermediate step
func BidRequestResponsePipe(ctx context.Context, bidRequestChan chan *BidRequest, bidResponseChan chan []*BidResponse, wg *sync.WaitGroup) {
	var x *BidRequest
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case x = <-bidRequestChan:
			bidResponseChan <- generateBidResponses(x)
			data, _ := json.MarshalIndent(&x, "", "   ")
			fmt.Println("Received data on bidRequestChannel:\n" + string(data))
		}
	}
}

func BidResponseConsumer(ctx context.Context, bidResponseChan chan []*BidResponse, auctionResultChan chan AuctionResult, wg *sync.WaitGroup) {
	var resArray []*BidResponse

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case resArray = <-bidResponseChan:
			data, _ := json.MarshalIndent(&resArray, "", "  ")
			fmt.Printf("Received data on bidResponseChannel:\n%s\n", string(data))

			auctionResult := simulateAuction(resArray)
			if auctionResult == nil {
				// e.g. empty bid array when every seat declined to bid
				continue
			}

			// gracefully log and discard if nothing is listening (i.e., ws not open)
			select {
			case auctionResultChan <- *auctionResult:
				auctionResultData, _ := json.MarshalIndent(auctionResult, "", " ")
				fmt.Printf("(successfully sent) Auction Result: %s\n", auctionResultData)
			default:
				auctionResultData, _ := json.MarshalIndent(&auctionResult, "", " ")
				fmt.Printf("(default block; discarded) Auction Result: %s\n", auctionResultData)
			}
		}
	}
}
