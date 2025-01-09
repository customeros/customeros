package dto

type SaveTag struct {
	Name       *string `json:"name,omitempty"`
	ColorCode  *string `json:"colorCode,omitempty"`
	EntityType *string `json:"entityType,omitempty"`
}
