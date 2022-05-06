package book

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ISBN        string    `json:"isbn"`
	Pages       int       `json:"pages"`
	Authors     []Author  `json:"authors"`
	Publisher   Publisher `json:"pub"`
	Description string    `json:"description"`
	Topics      []Topic   `json:"topics"`
	Language    Language  `json:"language"`
	Year        string    `json:"year"`
	Cover       Cover     `json:"cover"`

	Digital struct {
		Price    int   `json:"price"`
		Discount int   `json:"discount"` // Percentage
		PDF      EBook `json:"pdf"`
		EPUB     EBook `json:"epub"`
		DJVU     EBook `json:"djvu"`
		AZW      EBook `json:"azw"`
		TXT      EBook `json:"txt"`
		DOCX     EBook `json:"docx"`
	} `json:"digital"`

	Physical struct {
		Price    int `json:"price"`
		Discount int `json:"discount"`
		Stock    int `json:"stock"`
	} `json:"physical"`
}
