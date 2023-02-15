package dto

type TenantSettingsJiraDTO struct {
	JiraAPIToken *string `json:"jiraAPIToken"`
	JiraDomain   *string `json:"jiraDomain"`
	JiraEmail    *string `json:"jiraEmail"`
}
