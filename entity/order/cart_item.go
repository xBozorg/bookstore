package order

type CartItem struct {
	ID       uint     `json:"id"`
	BookId   uint     `json:"bookID"`
	UserID   string   `json:"userID"`
	Price    uint     `json:"price"`
	Date     string   `json:"timestamp"`
	Quantity uint     `json:"quantity"`
	Type     struct { // Digital | Physical | Bundle
		Digital  bool `json:"digital"`
		Physical bool `json:"physical"`
	} `json:"type"`
}