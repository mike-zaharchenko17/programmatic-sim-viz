package generator

func simulateAuction(bidResponseArray []*BidResponse) Bid {
	var max float64 = 0
	var maxBid Bid

	for _, res := range bidResponseArray {
		bid := res.SeatBid[0].Bid[0]
		if bid.Price > max {
			max = bid.Price
			maxBid = bid
		}
	}

	return maxBid
}
