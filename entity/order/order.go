package order

const (
	StatusCreated  uint = 100
	StatusPaid     uint = 110
	StatusVerified uint = 120
	StatusShipped  uint = 200
)

type Order struct {
	ID             uint   `json:"id"`
	UserID         string `json:"userID"`
	Total          uint   `json:"total"`
	Status         uint   `json:"status"`
	STN            string `json:"stn"` // Shipment Tracking Number
	CreationDate   string `json:"creationDate"`
	ReceiptionDate string `json:"receiptionDate"`
	Items          []Item `json:"Items"`
	Promo          Promo  `json:"promo"`
}
