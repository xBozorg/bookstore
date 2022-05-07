package user

type User struct {
	ID           string        `json:"id"` //UUID
	Email        string        `json:"email"`
	Password     string        `json:"password"`
	Username     string        `json:"username"`
	FirstName    string        `json:"firstName"`
	LastName     string        `json:"lastName"`
	PhoneNumbers []PhoneNumber `json:"phoneNumbers"`
	Addresses    []Address     `json:"addresses"`
	RegDate      string        `json:"regDate"` //Registration Date
}
