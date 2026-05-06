package generator

import "time"

// BidRequest is the top-level object for an auction request
type BidRequest struct {
	ID     string  `json:"id"`
	Imp    []Imp   `json:"imp"`
	Site   *Site   `json:"site,omitempty"`
	App    *App    `json:"app,omitempty"`
	Device *Device `json:"device,omitempty"`
	User   *User   `json:"user,omitempty"`
	At     int     `json:"at,omitempty"`
	Tmax   *int    `json:"tmax,omitempty"`
}

// Imp describes an impression being auctioned
type Imp struct {
	ID       string  `json:"id"`
	Bidfloor float64 `json:"bidfloor,omitempty"`
	TagID    *string `json:"tagid,omitempty"`
	Secure   *int    `json:"secure,omitempty"`
	Banner   *Banner `json:"banner,omitempty"`
	Video    *Video  `json:"video,omitempty"`
}

// Banner describes a display ad opportunity
type Banner struct {
	W   *int `json:"w,omitempty"`
	H   *int `json:"h,omitempty"`
	Pos *int `json:"pos,omitempty"`
}

// Video describes a video ad opportunity
type Video struct {
	W           *int `json:"w,omitempty"`
	H           *int `json:"h,omitempty"`
	Minduration *int `json:"minduration,omitempty"`
	Maxduration *int `json:"maxduration,omitempty"`
	Placement   *int `json:"placement,omitempty"`
}

// Site describes the publisher's website
type Site struct {
	ID        string     `json:"id,omitempty"`
	Domain    string     `json:"domain,omitempty"`
	Cat       []string   `json:"cat,omitempty"`
	Publisher *Publisher `json:"publisher,omitempty"`
}

// App describes the publisher's application
type App struct {
	ID        string     `json:"id,omitempty"`
	Bundle    string     `json:"bundle,omitempty"`
	Cat       []string   `json:"cat,omitempty"`
	Publisher *Publisher `json:"publisher,omitempty"`
}

// Publisher describes the seller
type Publisher struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Device describes the user's device
type Device struct {
	DeviceType     *int   `json:"devicetype,omitempty"`
	OS             string `json:"os,omitempty"`
	IP             string `json:"ip,omitempty"`
	IFA            string `json:"ifa,omitempty"`
	Geo            *Geo   `json:"geo,omitempty"`
	ConnectionType *int   `json:"connectiontype,omitempty"`
}

// Geo describes a geographic location
type Geo struct {
	Country string `json:"country,omitempty"`
	Region  string `json:"region,omitempty"`
}

// User describes the human user of the device
type User struct {
	ID  string `json:"id,omitempty"`
	Geo *Geo   `json:"geo,omitempty"`
}

// --- Bid Response Types ---

// BidResponse is the top-level response to a BidRequest
type BidResponse struct {
	ID      string    `json:"id"`
	SeatBid []SeatBid `json:"seatbid,omitempty"`
	Cur     string    `json:"cur,omitempty"`
	Nbr     *int      `json:"nbr,omitempty"`
}

// SeatBid contains bids from a specific buyer seat
type SeatBid struct {
	Bid  []Bid  `json:"bid"`
	Seat string `json:"seat,omitempty"`
}

// Bid represents an offer to buy an impression
type Bid struct {
	ID           string   `json:"id"`
	ImpID        string   `json:"impid"`
	Price        float64  `json:"price"`
	Adomain      []string `json:"adomain,omitempty"`
	CID          string   `json:"cid,omitempty"`
	CrID         string   `json:"crid,omitempty"`
	W            *int     `json:"w,omitempty"`
	H            *int     `json:"h,omitempty"`
	isBelowFloor bool
}

type LossRecord struct {
	Bid        *Bid `json:"bid"`
	LossReason int  `json:"loss_reason"`
}

type AuctionResult struct {
	RequestID     string       `json:"request_id"`
	Timestamp     time.Time    `json:"timestamp"`
	Winner        *Bid         `json:"winner,omitempty"`
	ClearingPrice float64      `json:"clearing_price"`
	Losers        []LossRecord `json:"losers,omitempty"`
}
