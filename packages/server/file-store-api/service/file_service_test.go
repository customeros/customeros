package service

import (
	"context"
	graph_model "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/test"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/config"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetById(t *testing.T) {
	// Setup test web server and mock resolver
	_, client, resolver := test.NewGraphQlMockedServer(t)

	resolver.Attachment = func(ctx context.Context, id string) (*graph_model.Attachment, error) {
		assert.Equal(t, "test-id", id)
		return &graph_model.Attachment{
			ID:        "test-id",
			MimeType:  "image/png",
			FileName:  "test.png",
			Size:      1024,
			BasePath:  "/GLOBAL",
			CdnURL:    "https://cdn.openline.com/test.png",
			CreatedAt: time.Now(),
		}, nil
	}

	// Create a file service with the mocked Config and GraphQL client
	fileSvc := &fileService{
		cfg:           &config.Config{},
		graphqlClient: client,
	}

	// Call GetById method
	file, err := fileSvc.GetById(context.Background(), "test-user", "test-tenant", "test-id")
	if err != nil {
		t.Fatalf("GetById returned error: %v", err)
	}

	assert.Equal(t, "test-id", file.ID)
	assert.Equal(t, "image/png", file.MimeType)
	assert.Equal(t, "test.png", file.FileName)
	assert.Equal(t, int64(1024), file.Size)
	assert.Equal(t, "/GLOBAL", file.BasePath)
	assert.Equal(t, "https://cdn.openline.com/test.png", file.CdnUrl)
}
