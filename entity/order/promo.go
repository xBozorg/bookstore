package order

type Promo struct {
	ID    uint   `json:"id"`
	Code  string `json:"code"`
	Offer struct {
		Percentage uint `json:"percentage"`
		Amount     uint `json:"amount"`
	} `json:"offer"`
}
