package user

type PhoneNumber struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Number string `json:"number"`
}
