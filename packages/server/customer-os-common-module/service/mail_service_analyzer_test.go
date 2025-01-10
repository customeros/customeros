package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAutoResponder(t *testing.T) {
	svc := &mailService{}
	tests := []struct {
		name     string
		headers  EmailHeaders
		wantBool bool
		wantMsg  string
	}{
		{
			name:     "X-Autoreply present",
			headers:  EmailHeaders{XAutoreply: "yes"},
			wantBool: true,
			wantMsg:  "AUTORESPONDER | X-AUTOREPLY",
		},
		{
			name:     "X-Autoreply empty string",
			headers:  EmailHeaders{XAutoreply: ""},
			wantBool: false,
			wantMsg:  "",
		},
		{
			name:     "X-Autoresponse present",
			headers:  EmailHeaders{XAutoresponse: "yes"},
			wantBool: true,
			wantMsg:  "AUTORESPONDER | X-AUTORESPONSE",
		},
		{
			name:     "X-Loop present",
			headers:  EmailHeaders{XLoop: true},
			wantBool: true,
			wantMsg:  "AUTORESPONDER | X-LOOP",
		},
		{
			name:     "Precedence auto_reply",
			headers:  EmailHeaders{Precedence: "auto_reply"},
			wantBool: true,
			wantMsg:  "AUTORESPONDER | PRECEDENCE: AUTO_REPLY",
		},
		{
			name:     "No autoresponder headers",
			headers:  EmailHeaders{},
			wantBool: false,
			wantMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBool, gotMsg := svc.isAutoResponder(tt.headers)
			assert.Equal(t, tt.wantBool, gotBool)
			assert.Equal(t, tt.wantMsg, gotMsg)
		})
	}
}

func TestIsBounce(t *testing.T) {
	svc := &mailService{}
	tests := []struct {
		name     string
		headers  EmailHeaders
		subject  string
		from     string
		wantBool bool
		wantMsg  string
	}{
		{
			name: "X-Failed-Recipients with emails",
			headers: EmailHeaders{
				XFailedRecepients: []string{"failed@example.com"},
			},
			wantBool: true,
			wantMsg:  "BOUNCE | X-FAILED-RECIPIENTS",
		},
		{
			name:     "Content-Description delivery report",
			headers:  EmailHeaders{ContentDescription: "delivery report"},
			wantBool: true,
			wantMsg:  "BOUNCE | CONTENT-DESCRIPTION: DELIVERY REPORT",
		},
		{
			name:     "Return-Path contains mailer-daemon",
			headers:  EmailHeaders{ReturnPath: "mailer-daemon@example.com"},
			wantBool: true,
			wantMsg:  "BOUNCE | RETURN-PATH CONTAINS BOUNCE KEYWORDS",
		},
		{
			name:     "From contains MAILER-DAEMON case insensitive",
			from:     "MAILER-DAEMON@example.com",
			wantBool: true,
			wantMsg:  "BOUNCE | FROM CONTAINS BOUNCE KEYWORDS",
		},
		{
			name:     "Not a bounce",
			headers:  EmailHeaders{},
			subject:  "Regular subject",
			from:     "user@example.com",
			wantBool: false,
			wantMsg:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBool, gotMsg := svc.isBounce(tt.headers, tt.subject, tt.from)
			assert.Equal(t, tt.wantBool, gotBool)
			assert.Equal(t, tt.wantMsg, gotMsg)
		})
	}
}

func TestIsBulkMail(t *testing.T) {
	svc := &mailService{}
	tests := []struct {
		name     string
		headers  EmailHeaders
		from     string
		replyTo  []EmailParticipant
		wantBool bool
		wantMsg  string
	}{
		{
			name: "Different Reply-To",
			headers: EmailHeaders{
				ReplyToExists: true,
				ForwardedFor:  "",
			},
			from: "sender@example.com",
			replyTo: []EmailParticipant{
				{Email: "different@example.com"},
			},
			wantBool: true,
			wantMsg:  "BULK | REPLY-TO != FROM",
		},
		{
			name: "Different Reply-To",
			headers: EmailHeaders{
				ReplyToExists: true,
				ForwardedFor:  "xxx",
			},
			from: "sender@example.com",
			replyTo: []EmailParticipant{
				{Email: "different@example.com"},
			},
			wantBool: false,
			wantMsg:  "",
		},
		{
			name: "Multiple Reply-To with match",
			headers: EmailHeaders{
				ReplyToExists: true,
				ReturnPath:    "sender@example.com",
			},
			from: "sender@example.com",
			replyTo: []EmailParticipant{
				{Email: "different1@example.com"},
				{Email: "sender@example.com"},
				{Email: "different2@example.com"},
			},
			wantBool: false,
			wantMsg:  "",
		},
		{
			name: "Empty Return-Path with flag",
			headers: EmailHeaders{
				ReturnPathExists: true,
				ForwardedFor:     "",
				ReturnPath:       "",
			},
			from: "sender@example.com",
			replyTo: []EmailParticipant{
				{Email: "sender@example.com"},
			},
			wantBool: true,
			wantMsg:  "BULK | EMPTY RETURN-PATH",
		},
		{
			name: "Different Return-Path with flag",
			headers: EmailHeaders{
				ReturnPathExists: true,
				ForwardedFor:     "",
				ReturnPath:       "different@example.com",
			},
			from: "sender@example.com",
			replyTo: []EmailParticipant{
				{Email: "sender@example.com"},
			},
			wantBool: true,
			wantMsg:  "BULK | RETURN-PATH != FROM",
		},
		{
			name: "Different Sender",
			headers: EmailHeaders{
				ReturnPathExists: true,
				ReturnPath:       "sender@example.com",
				Sender:           "different@example.com",
			},
			from: "sender@example.com",
			replyTo: []EmailParticipant{
				{Email: "sender@example.com"},
			},
			wantBool: true,
			wantMsg:  "BULK | SENDER != FROM",
		},
		{
			name: "List-Unsubscribe present",
			headers: EmailHeaders{
				ListUnsubscribe: true,
			},
			from: "sender@example.com",
			replyTo: []EmailParticipant{
				{Email: "sender@example.com"},
			},
			wantBool: true,
			wantMsg:  "BULK | UNSUBSCRIBE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBool, gotMsg := svc.isBulkMail(tt.headers, tt.from, tt.replyTo)
			assert.Equal(t, tt.wantBool, gotBool, "test: %s", tt.name)
			assert.Equal(t, tt.wantMsg, gotMsg, "test: %s", tt.name)
		})
	}
}

func TestIsReturnPathBounce(t *testing.T) {
	svc := &mailService{}
	tests := []struct {
		name       string
		returnPath string
		want       bool
	}{
		{
			name:       "Contains mailer-daemon",
			returnPath: "mailer-daemon@example.com",
			want:       true,
		},
		{
			name:       "Contains MAILER-DAEMON case insensitive",
			returnPath: "MAILER-DAEMON@example.com",
			want:       true,
		},
		{
			name:       "Normal return path",
			returnPath: "user@example.com",
			want:       false,
		},
		{
			name:       "Empty return path",
			returnPath: "",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.isReturnPathBounce(tt.returnPath)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsBounceSubject(t *testing.T) {
	svc := &mailService{}
	tests := []struct {
		name    string
		subject string
		want    bool
	}{
		{
			name:    "Delivery Status Notification",
			subject: "Delivery Status Notification (Failure)",
			want:    true,
		},
		{
			name:    "DELIVERY STATUS NOTIFICATION case insensitive",
			subject: "DELIVERY STATUS NOTIFICATION (FAILURE)",
			want:    true,
		},
		{
			name:    "Undeliverable",
			subject: "Undeliverable: Your message",
			want:    true,
		},
		{
			name:    "Undelivered",
			subject: "Undelivered Mail Returned",
			want:    true,
		},
		{
			name:    "Delivery Failure",
			subject: "Delivery Failure Notice",
			want:    true,
		},
		{
			name:    "Regular Subject",
			subject: "Meeting Tomorrow",
			want:    false,
		},
		{
			name:    "Empty Subject",
			subject: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.isBounceSubject(tt.subject)
			assert.Equal(t, tt.want, got)
		})
	}
}
