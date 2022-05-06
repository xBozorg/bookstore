package user

type PhoneNumber struct {
	ID     uint   `json:"id"`
	Code   string `json:"code"`
	Number string `json:"number"`
}
