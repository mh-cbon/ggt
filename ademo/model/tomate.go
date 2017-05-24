package model

type Tomate struct {
	ID    string
	Color string
}

func (t Tomate) GetID() string {
	return t.ID
}
