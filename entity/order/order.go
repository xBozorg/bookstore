package order

type Order struct {
	ID         uint   `json:"id"`
	UserID     string `json:"userID"`
	CartID     uint   `json:"cartID"`
	Total      uint   `json:"total"`
	Date       string `json:"date"`
	StatusCode uint   `json:"statusCode"`
	STN        uint   `json:"stn"` // Shipment Tracking Number
}
