package dto

type AddTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewAddTagEvent(id, name string) AddTag {
	output := AddTag{
		ID:   id,
		Name: name,
	}
	return output
}
