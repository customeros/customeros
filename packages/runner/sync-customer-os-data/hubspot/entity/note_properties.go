package entity

type NoteProperties struct {
	AirbyteAbId        string `gorm:"column:_airbyte_ab_id"`
	AirbyteNotesHashid string `gorm:"column:_airbyte_engagements_notes_hashid"`
	NoteBody           string `gorm:"column:hs_note_body"`
	OwnerId            string `gorm:"column:hubspot_owner_id"`
}

type NotePropertiesList []NoteProperties

func (NoteProperties) TableName() string {
	return "engagements_notes_properties"
}
