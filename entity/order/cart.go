package order

type Cart struct {
	ID     uint       `json:"ID"`
	UserID string     `json:"userID"`
	Total  uint       `json:"price"`
	Items  []CartItem `json:"orders"`
	Promo  Promo      `json:"promo"`
}
