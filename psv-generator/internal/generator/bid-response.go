package generator

import "math/rand"

func generateBidResponses(bidRequest *BidRequest) []*BidResponse {

	var responses []*BidResponse

	// every seat has a chance to bid
	for _, seat := range Seats {
		// 0.4 pr of not bidding
		if rand.Float32() > 0.6 {
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
