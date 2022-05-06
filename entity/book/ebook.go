package book

type EBook struct {
	Available bool   `json:"available"`
	Path      string `json:"path"`
}
