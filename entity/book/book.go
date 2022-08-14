package book

const (
	NotAvailable uint = iota
	DigitalAvailable
	PhysicalAvailable
	BundleAvailable
)

type Book struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	ISBN         string    `json:"isbn"`
	Pages        uint      `json:"pages"`
	Authors      []Author  `json:"authors"`
	Publisher    Publisher `json:"pub"`
	Description  string    `json:"description"`
	Topics       []Topic   `json:"topics"`
	Language     Language  `json:"language"`
	Year         string    `json:"year"`
	CoverFront   string    `json:"coverFront"`
	CoverBack    string    `json:"coverBack"`
	CreationDate string    `json:"creationDate"`
	Digital      Digital   `json:"digital"`
	Physical     Physical  `json:"physical"`
	Availability uint      `json:"availability"`
}

type Digital struct {
	Price    uint   `json:"price"`
	Discount uint   `json:"discount"` // Percentage
	PDF      string `json:"pdf"`
	EPUB     string `json:"epub"`
	DJVU     string `json:"djvu"`
	AZW      string `json:"azw"`
	TXT      string `json:"txt"`
	DOCX     string `json:"docx"`
}

type Physical struct {
	Price    uint `json:"price"`
	Discount uint `json:"discount"`
	Stock    uint `json:"stock"`
}
