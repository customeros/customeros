package dto

type RemoveTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewRemoveTagEvent(id, name string) RemoveTag {
	output := RemoveTag{
		ID:   id,
		Name: name,
	}
	return output
}
