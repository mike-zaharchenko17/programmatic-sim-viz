package generator

import (
	"math/rand"
	"psv-generator/internal/utils"
)

func GenerateBidRequest() *BidRequest {
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
			Banner:   &Banner{W: utils.ToPtr(300), H: utils.ToPtr(300)},
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
