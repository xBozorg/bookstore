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
	PhoneID        uint   `json:"phoneID"`
	AddressID      uint   `json:"addressID"`
}

type OrderPaymentInfo struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	Total uint   `json:"total"`
}

type ZarinpalOrder struct {
	ID        uint   `json:"id"`
	OrderID   uint   `json:"orderID"`
	Authority string `json:"authority"`
	RefID     int    `json:"refID"`
	Code      int    `json:"code"`
}
