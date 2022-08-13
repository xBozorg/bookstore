package order

type Promo struct {
	ID         uint   `json:"id"`
	Code       string `json:"code"`
	Percentage uint   `json:"percentage"`
	Expiration string `json:"expiration"`
	Limit      uint   `json:"limit"`
	MaxPrice   uint   `json:"maxPrice"`
}
