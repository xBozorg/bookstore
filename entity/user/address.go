package user

type Address struct {
	ID          uint   `json:"id"`
	Country     string `json:"country"`
	Province    string `json:"province"`
	City        string `json:"city"`
	Street      string `json:"street"`
	PostalCode  string `json:"postalCode"`
	No          uint   `json:"no"`
	Description string `json:"description"`
}
