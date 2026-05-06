package generator

import (
	"math/rand"
	"time"
)

func simulateAuction(bidResponseArray []*BidResponse) *AuctionResult {
	if len(bidResponseArray) == 0 {
		return nil
	}

	// roll opportunity expired- 5% chance; reject all
	if rand.Intn(100) <= 5 {
		auctionResult := AuctionResult{
			RequestID:     bidResponseArray[0].ID,
			Timestamp:     time.Now(),
			Winner:        nil,
			ClearingPrice: 0,
		}

		for _, res := range bidResponseArray {
			bid := res.SeatBid[0].Bid[0]
			auctionResult.Losers = append(auctionResult.Losers, LossRecord{
				Bid:        &bid,
				LossReason: 2,
			})
		}

		return &auctionResult
	}

	var maxBid Bid = bidResponseArray[0].SeatBid[0].Bid[0]
	var maxPrice = maxBid.Price

	for _, res := range bidResponseArray {
		bid := res.SeatBid[0].Bid[0]

		if bid.Price > maxPrice {
			maxBid = bid
			maxPrice = bid.Price
		}
	}

	auctionResult := AuctionResult{
		RequestID:     bidResponseArray[0].ID,
		Timestamp:     time.Now(),
		Winner:        &maxBid,
		ClearingPrice: maxBid.Price,
	}

	for _, res := range bidResponseArray {
		bid := res.SeatBid[0].Bid[0]

		if maxBid.ID != bid.ID {
			if bid.isBelowFloor {
				auctionResult.Losers = append(auctionResult.Losers, LossRecord{
					Bid:        &bid,
					LossReason: 100,
				})
			} else {
				auctionResult.Losers = append(auctionResult.Losers, LossRecord{
					Bid:        &bid,
					LossReason: 102,
				})
			}
		}
	}

	return &auctionResult
}
