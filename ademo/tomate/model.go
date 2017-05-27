package tomate

// Tomate is a model of tomatoes
type Tomate struct {
	ID    string
	Color string
}

// GetID is useful for identity check.
func (t Tomate) GetID() string {
	return t.ID
}

// SimilarTomate indiicates tomate similarity to a value
type SimilarTomate struct {
	Tomate
	Similarity float64
}
