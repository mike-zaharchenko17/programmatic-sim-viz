package generator

import "time"

func simulateAuction(bidResponseArray []*BidResponse) AuctionResult {
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
			auctionResult.Losers = append(auctionResult.Losers, LossRecord{
				Bid:        &bid,
				LossReason: 102,
			})
		}
	}

	return auctionResult
}
