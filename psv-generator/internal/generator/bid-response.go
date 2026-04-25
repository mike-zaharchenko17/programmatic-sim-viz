package generator

import "math/rand"

func generateBidResponses(bidRequest *BidRequest) []*BidResponse {

	var responses []*BidResponse

	// every seat has a chance to bid
	for _, seat := range Seats {
		// 0.25 pr of not bidding
		// in a truer to life sim, we would return a BidResponse
		// with a nbr code, but instead we omit to keep things lean
		if rand.Float32() > 0.75 {
			continue
		}

		// randomly select one of the campaigns that belongs to the seat
		campaign := PickOne(seat.Campaigns)

		// select one of the impressions on the bid request
		imp := PickOne(bidRequest.Imp)
		bidPrice := GenBidPrice(imp.Bidfloor)

		responses = append(responses, &BidResponse{
			ID: bidRequest.ID,
			SeatBid: []SeatBid{{
				Seat: seat.Seat,
				Bid: []Bid{{
					ID:      GenUUID(),
					ImpID:   imp.ID,
					Price:   bidPrice,
					CID:     campaign.CID,
					CrID:    GenUUID(),
					Adomain: []string{campaign.Adomain},
				}},
			}},
		})
	}

	return responses
}
