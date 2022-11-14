package service

import (
	"fmt"
	"github.com/machinebox/graphql"
	"testing"
)

func Test_getContactByEmail(t *testing.T) {
	graphqlClient := graphql.NewClient("http://localhost:10010/query")
	contact, err := getContactByEmail(graphqlClient, "tsearle@gmail.com")

	if err != nil {
		t.Errorf("Got an error %v", err)
		return
	}
	fmt.Printf("Got a contact of %s %s %s", contact.firstName, contact.lastName, contact.id)
}

func Test_createContact(t *testing.T) {
	graphqlClient := graphql.NewClient("http://localhost:10010/query")

	contact, err := createContact(graphqlClient, "Torrey", "Searle", "tsearle@gmail.com")

	if err != nil {
		t.Errorf("Got an error %v", err)
		return
	}
	fmt.Printf("Got a contact of %s", contact)
}

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
