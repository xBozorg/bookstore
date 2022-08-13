package order

const (
	Digital uint = iota
	Physical
	Bundle
)

type Item struct {
	ID       uint `json:"id"`
	BookID   uint `json:"bookID"`
	Type     uint `json:"type"`
	Quantity uint `json:"quantity"`
}
