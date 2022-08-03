package admin

type Admin struct {
	ID          string `json:"id"` //UUID
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}
