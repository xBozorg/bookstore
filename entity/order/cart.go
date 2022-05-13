package order

type Cart struct {
	ID     uint       `json:"id"`
	UserID string     `json:"userID"`
	Total  uint       `json:"total"`
	Items  []CartItem `json:"items"`
	Promo  Promo      `json:"promo"`
}
