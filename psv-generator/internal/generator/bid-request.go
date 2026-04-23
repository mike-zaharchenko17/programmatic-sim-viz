package generator

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"psv-generator/internal/utils"
	"sync"
	"time"
)

func generateBidRequest() *BidRequest {
	isApp := rand.Float32() < 0.30

	var site *Site
	var app *App
	var publisherId string

	if isApp {
		a := PickOne(Apps)
		publisherId = a.PublisherID
		app = &App{ID: a.ID, Bundle: a.Bundle, Cat: a.Cat}
	} else {
		s := PickOne(Sites)
		publisherId = s.PublisherID
		site = &Site{ID: s.ID, Domain: s.Domain, Cat: s.Cat}
	}

	var pub *Publisher
	for _, p := range Publishers {
		if p.ID == publisherId {
			pub = &Publisher{ID: p.ID, Name: p.Name}
			break
		}
	}

	if site != nil {
		site.Publisher = pub
	} else {
		app.Publisher = pub
	}

	deviceType := PickWeighted(DeviceTypes)
	connType := PickWeighted(ConnectionTypes)

	return &BidRequest{
		ID: GenIFA(),
		Imp: []Imp{{
			ID:       "1",
			Bidfloor: PickOne(BidFloors),
			Banner:   &Banner{W: utils.IntPtr(300), H: utils.IntPtr(300)},
		}},
		Site: site,
		App:  app,
		Device: &Device{
			DeviceType:     &deviceType,
			OS:             PickWeighted(OSTypes),
			IP:             GenIP(),
			IFA:            GenIFA(),
			Geo:            &Geo{Country: PickWeighted(Countries)},
			ConnectionType: &connType,
		},
		User: &User{ID: GenIFA()},
		At:   2,
	}
}

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

func BidRequestConsumer(ctx context.Context, bidRequestChan chan *BidRequest, wg *sync.WaitGroup) {
	var x *BidRequest
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		case x = <-bidRequestChan:
			data, _ := json.MarshalIndent(&x, "", "   ")
			fmt.Println(string(data))
		}
	}
}
