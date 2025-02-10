package postgres_entity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/pkg/errors"
	"time"
)

type OAuthTokenEntity struct {
	Provider                  string    `gorm:"primaryKey;autoIncrement:false;index:idx_primary;column:provider;size:255;not null"`
	TenantName                string    `gorm:"primaryKey;autoIncrement:false;index:idx_primary;column:tenant_name;size:255;not null"`
	EmailAddress              string    `gorm:"primaryKey;autoIncrement:false;index:idx_primary;column:email_address;size:255;not null"`
	Type                      string    `gorm:"column:type;size:50;"`
	PlayerIdentityId          string    `gorm:"column:player_identity_id;size:255;not null"`
	AccessToken               string    `gorm:"column:access_token;type:text"`
	RefreshToken              string    `gorm:"column:refresh_token;type:text"`
	NeedsManualRefresh        bool      `gorm:"column:needs_manual_refresh;default:false;"`
	IdToken                   string    `gorm:"column:id_token;type:text"`
	ExpiresAt                 time.Time `gorm:"column:expires_at;type:timestamp;"`
	Scope                     string    `gorm:"column:scope;type:text"`
	GmailSyncEnabled          bool      `gorm:"column:gmail_sync_enabled;default:false;"`
	GoogleCalendarSyncEnabled bool      `gorm:"column:google_calendar_sync_enabled;default:false;"`
}

func (OAuthTokenEntity) TableName() string {
	return "oauth_token"
}

// Encrypt a token using AES-GCM
func EncryptToken(key string, token string) (string, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(decodedKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt a token using AES-GCM
func DecryptToken(key string, encryptedToken string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(decodedKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
