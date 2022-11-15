package service

import (
	"testing"
)

func Test_parseEmail(t *testing.T) {
	type args struct {
		email string
	}

	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{name: "unquoted displayname",
			args:  args{email: "Torrey Searle <tsearle@invalid.domain>"},
			want:  "Torrey Searle",
			want1: "tsearle@invalid.domain",
		},
		{name: "quoted displayname",
			args:  args{email: "\"Torrey Searle\" <tsearle@invalid.domain>"},
			want:  "Torrey Searle",
			want1: "tsearle@invalid.domain",
		},
		{name: "no display name",
			args:  args{email: "tsearle@invalid.domain"},
			want:  "",
			want1: "tsearle@invalid.domain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := parseEmail(tt.args.email)
			if got != tt.want {
				t.Errorf("parseEmail() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("parseEmail() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
