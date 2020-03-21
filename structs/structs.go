package structs

type User struct {
	Id    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

type Bid struct {
	BidId            string  `json:"bid_id" db:"bid_id"`
	UserId           string  `json:"user_id" db:"user_id"`
	ItemId           string  `json:"item_id" db:"item_id"`
	BidAmount        float32 `json:"bid_amount" db:"bid_amount"`
	IsAvailableToBid bool    `json:"is_available_to_bid" db:"is_available_to_bid"`
}

type Item struct {
	Id    string  `json:"id" db:"id"`
	Name  string  `json:"name" db:"name"`
	Price float32 `json:"price" db:"price"`
}
