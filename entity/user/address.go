package user

type Address struct {
	ID          int    `json:"id"`
	Country     string `json:"country"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Street      string `json:"street"`
	PostalCode  string `json:"postalCode"`
	No          int    `json:"no"`
	Description string `json:"description"`
}
