package resolver

import (
	"context"
	"github.com/99designs/gqlgen/client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils/decode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryResolver_TableViewDefs(t *testing.T) {
	ctx := context.TODO()
	defer tearDownTestCase(ctx)(t)

	rawResponse, err := c.RawPost(getQuery("view/get_views"),
		client.Var("page", 1),
		client.Var("limit", 1),
	)
	assertRawResponseSuccess(t, rawResponse, err)

	var tableViewDefs struct {
		TableViewDefs model.TableViewDefPage
	}

	err = decode.Decode(rawResponse.Data.(map[string]any), &tableViewDefs)
	require.Nil(t, err)
	require.NotNil(t, tableViewDefs)
	pagedTableViewDefs := tableViewDefs.TableViewDefs
	require.Equal(t, 1, pagedTableViewDefs.TotalPages)
	require.Equal(t, int64(1), pagedTableViewDefs.TotalElements)
	require.Equal(t, model.TableViewDef{ID: "MockID", Name: "MockName"}, *pagedTableViewDefs.Content[0])
}
