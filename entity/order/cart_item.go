package order

type CartItem struct {
	ID       uint     `json:"id"`
	BookID   uint     `json:"bookID"`
	UserID   string   `json:"userID"`
	Price    uint     `json:"price"`
	Date     string   `json:"date"`
	Quantity uint     `json:"quantity"`
	Type     struct { // Digital | Physical | Bundle
		Digital  bool `json:"digital"`
		Physical bool `json:"physical"`
	} `json:"type"`
}
