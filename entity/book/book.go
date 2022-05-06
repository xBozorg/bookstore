package book

type Book struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	ISBN        string    `json:"isbn"`
	Pages       uint      `json:"pages"`
	Authors     []Author  `json:"authors"`
	Publisher   Publisher `json:"pub"`
	Description string    `json:"description"`
	Topics      []Topic   `json:"topics"`
	Language    Language  `json:"language"`
	Year        string    `json:"year"`
	Cover       Cover     `json:"cover"`

	Digital struct {
		Price    uint  `json:"price"`
		Discount uint  `json:"discount"` // Percentage
		PDF      EBook `json:"pdf"`
		EPUB     EBook `json:"epub"`
		DJVU     EBook `json:"djvu"`
		AZW      EBook `json:"azw"`
		TXT      EBook `json:"txt"`
		DOCX     EBook `json:"docx"`
	} `json:"digital"`

	Physical struct {
		Price    uint `json:"price"`
		Discount uint `json:"discount"`
		Stock    uint `json:"stock"`
	} `json:"physical"`
}
