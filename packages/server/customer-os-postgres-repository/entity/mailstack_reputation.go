package postgres_entity

import "time"

type MailstackReputationEntity struct {
	ID                  string    `gorm:"primary_key;type:uuid;default:gen_random_uuid()" json:"id"`
	Tenant              string    `gorm:"column:tenant;type:varchar(255);NOT NULL" json:"tenant"`
	CreatedAt           time.Time `gorm:"column:created_at;type:timestamp;DEFAULT:current_timestamp" json:"createdAt"`
	Domain              string    `gorm:"column:domain;type:varchar(255)" json:"domain"`
	DomainAgePenalty    int       `gorm:"column:domain_age_penalty;type:integer" json:"domainAgePenalty"`
	BlacklistPenaltyPct int       `gorm:"column:blacklist_penalty_pct;type:integer" json:"blacklistPenaltyPct"`
	BouncePenaltyPct    int       `gorm:"column:bounce_penalty_pct;type:integer" json:"bouncePenaltyPct"`
	DMARCPenaltyPct     int       `gorm:"column:dmarc_penalty_pct;type:integer" json:"dmarcPenaltyPct"`
	SPFPenaltyPct       int       `gorm:"column:spf_penalty_pct;type:integer" json:"spfPenaltyPct"`
}

func (MailstackReputationEntity) TableName() string {
	return "mailstack_reputation"
}
