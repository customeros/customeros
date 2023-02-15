package dto

type TenantSettingsTrelloDTO struct {
	TrelloAPIToken *string `json:"trelloAPIToken"`
	TrelloAPIKey   *string `json:"trelloAPIKey"`
}
