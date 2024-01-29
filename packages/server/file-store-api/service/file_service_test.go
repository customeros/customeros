package service

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/model"
	graph_model "github.com/openline-ai/openline-customer-os/packages/server/file-store-api/test/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/test/utils"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestGetById(t *testing.T) {
	// Setup test web server and mock resolver
	_, client, resolver := utils.NewWebServer(t)

	resolver.Attachment = func(ctx context.Context, id string) (*graph_model.Attachment, error) {
		assert.Equal(t, "test-id", id)

		return &graph_model.Attachment{
			ID:            "test-id",
			CreatedAt:     time.Now().UTC(),
			MimeType:      "image/png",
			FileName:      "test.png",
			Size:          1024,
			BasePath:      "/GLOBAL",
			Source:        "TEST",
			SourceOfTruth: "TEST",
			AppSource:     "file-store-api",
		}, nil
	}

	// Create a file service with the mocked Config and GraphQL client
	fileSvc := &fileService{
		cfg:           &config.Config{},
		graphqlClient: client,
	}

	// Call GetById method
	file, err := fileSvc.GetById("test-user", "test-tenant", "test-id")
	if err != nil {
		t.Fatalf("GetById returned error: %v", err)
	}

	// Verify that the file object matches the expected value
	expectedFile := &model.File{
		ID:       "test-id",
		MimeType: "image/png",
		FileName: "test.png",
		Size:     1024,
		BasePath: "/GLOBAL",
	}
	if !reflect.DeepEqual(file, expectedFile) {
		t.Fatalf("Unexpected file object: %v", file)
	}
}
