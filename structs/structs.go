package structs

type User struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type Bid struct {
	BidId            string  `json:"bid_id"`
	UserId           string  `json:"user_id"`
	ItemId           string  `json:"item_id"`
	BidAmount        float32 `json:"bid_amount"`
}

type Item struct {
	ItemId string `json:"item_id"`
	Name   string  `json:"name"`
	Price  float32 `json:"price"`
	IsAvailableToBid bool `json:"is_available_to_bid"`
}
