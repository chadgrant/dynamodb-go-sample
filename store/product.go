package store

type Product struct {
	ID          string  `json:"id"`
	Category    string  `json:"category"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description,omitempty"`
}
