package order

type Order struct {
	ID         uint   `json:"id"`
	UserID     string `json:"userID"`
	CartId     uint   `json:"cartId"`
	Total      uint   `json:"total"`
	Date       string `json:"timestamp"`
	StatusCode uint   `json:"statusCode"`
	STN        uint   `json:"stn"` // Shipment Tracking Number
}
