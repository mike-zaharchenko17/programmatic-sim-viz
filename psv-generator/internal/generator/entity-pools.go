package generator

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
)

var CampaignsBySeat map[string][]CampaignPool
var Seats []SeatPool

func init() {
	CampaignsBySeat = make(map[string][]CampaignPool)
	for _, c := range Campaigns {
		CampaignsBySeat[c.Seat] = append(CampaignsBySeat[c.Seat], c)
	}

	for seat, campaigns := range CampaignsBySeat {
		Seats = append(Seats, SeatPool{
			Seat:      seat,
			Campaigns: campaigns,
		})
	}
}

// Entity Pools
type PublisherPool struct {
	ID   string
	Name string
}
type SitePool struct {
	ID          string
	Domain      string
	PublisherID string
	Cat         []string
}
type AppPool struct {
	ID          string
	Bundle      string
	PublisherID string
	Cat         []string
}
type CampaignPool struct {
	CID     string
	Seat    string
	Adomain string
}

type SeatPool struct {
	Seat      string
	Campaigns []CampaignPool
}

// Pool Data
var Publishers = []PublisherPool{
	{ID: "pub-001", Name: "NewsMedia Inc"},
	{ID: "pub-002", Name: "SportsCo"},
	{ID: "pub-003", Name: "GamerNetwork"},
	{ID: "pub-004", Name: "LifestyleDigital"},
	{ID: "pub-005", Name: "TechDaily"},
}
var Sites = []SitePool{
	{ID: "site-001", Domain: "news-example.com", PublisherID: "pub-001", Cat: []string{"IAB12"}},
	{ID: "site-002", Domain: "sports-example.com", PublisherID: "pub-002", Cat: []string{"IAB17"}},
	{ID: "site-003", Domain: "gaming-example.com", PublisherID: "pub-003", Cat: []string{"IAB9"}},
	{ID: "site-004", Domain: "lifestyle-example.com", PublisherID: "pub-004", Cat: []string{"IAB18"}},
	{ID: "site-005", Domain: "tech-example.com", PublisherID: "pub-005", Cat: []string{"IAB19"}},
}
var Apps = []AppPool{
	{ID: "app-001", Bundle: "com.newsmedia.app", PublisherID: "pub-001", Cat: []string{"IAB12"}},
	{ID: "app-002", Bundle: "com.sportsco.live", PublisherID: "pub-002", Cat: []string{"IAB17"}},
	{ID: "app-003", Bundle: "com.gamernetwork.hub", PublisherID: "pub-003", Cat: []string{"IAB9"}},
}
var Campaigns = []CampaignPool{
	{CID: "camp-001", Seat: "seat-acme", Adomain: "acme-brand.com"},
	{CID: "camp-002", Seat: "seat-acme", Adomain: "acme-brand.com"},
	{CID: "camp-003", Seat: "seat-bigretail", Adomain: "bigretail.com"},
	{CID: "camp-004", Seat: "seat-bigretail", Adomain: "bigretail.com"},
	{CID: "camp-005", Seat: "seat-automaker", Adomain: "automaker.com"},
	{CID: "camp-006", Seat: "seat-streamco", Adomain: "streamco.com"},
	{CID: "camp-007", Seat: "seat-finserv", Adomain: "finserv.com"},
	{CID: "camp-008", Seat: "seat-finserv", Adomain: "finserv.com"},
	{CID: "camp-009", Seat: "seat-bigretail", Adomain: "bigretail.com"},
}

// Common bid floors
var BidFloors = []float64{0.10, 0.25, 0.50, 1.00, 2.00, 5.00}

type Weighted[V int | string] struct {
	Value  V
	Weight int
}

// Device type weights (OpenRTB 5.21: 1=mobile, 2=PC, 3=connected TV, 4=phone, 5=tablet)
var DeviceTypes = []Weighted[int]{
	{Value: 4, Weight: 40}, // phone
	{Value: 5, Weight: 20}, // tablet
	{Value: 2, Weight: 25}, // PC
	{Value: 3, Weight: 10}, // CTV
	{Value: 1, Weight: 5},  // mobile (generic)
}
var OSTypes = []Weighted[string]{
	{Value: "iOS", Weight: 35},
	{Value: "Android", Weight: 40},
	{Value: "Windows", Weight: 15},
	{Value: "macOS", Weight: 8},
	{Value: "Other", Weight: 2},
}
var ConnectionTypes = []Weighted[int]{
	{Value: 2, Weight: 50}, // wifi
	{Value: 6, Weight: 30}, // cellular 4g
	{Value: 7, Weight: 10}, // cellular 5g
	{Value: 1, Weight: 5},  // ethernet
	{Value: 0, Weight: 5},  // unknown
}
var Countries = []Weighted[string]{
	{Value: "USA", Weight: 50},
	{Value: "GBR", Weight: 12},
	{Value: "DEU", Weight: 8},
	{Value: "CAN", Weight: 7},
	{Value: "AUS", Weight: 6},
	{Value: "FRA", Weight: 5},
	{Value: "JPN", Weight: 5},
	{Value: "BRA", Weight: 4},
	{Value: "IND", Weight: 3},
}

func PickWeighted[V int | string](choices []Weighted[V]) V {
	total := 0
	for _, c := range choices {
		total += c.Weight
	}
	r := rand.Intn(total)
	for _, c := range choices {
		r -= c.Weight
		if r < 0 {
			return c.Value
		}
	}

	return choices[0].Value
}

func PickOne[T any](items []T) T {
	idx := rand.Intn(len(items))
	return items[idx]
}

// log normal-ish bid price above floor
func GenBidPrice(floor float64) float64 {
	spread := floor * 0.5
	if spread < 0.10 {
		spread = 0.10
	}
	return floor + rand.ExpFloat64()*spread
}

func GenIP() string {
	s := fmt.Sprintf("%d.%d.%d.%d",
		rand.Intn(224), rand.Intn(256), rand.Intn(256), rand.Intn(256))

	return s
}

func GenIFA() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		rand.Uint32(), rand.Intn(0xFFFF), rand.Intn(0xFFFF),
		rand.Intn(0xFFFF), rand.Int63n(0xFFFFFFFFFFFF))
}

func GenUUID() string {
	uuidAsString := uuid.NewString()
	return uuidAsString
}
