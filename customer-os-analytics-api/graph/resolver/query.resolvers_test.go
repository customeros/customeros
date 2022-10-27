package resolver

import (
	"encoding/json"
	"github.com.openline-ai.customer-os-analytics-api/graph/generated"
	"github.com.openline-ai.customer-os-analytics-api/mocks"
	"github.com.openline-ai.customer-os-analytics-api/repository"
	"github.com.openline-ai.customer-os-analytics-api/repository/helper"
	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestQueryApplicationById_NotFound(t *testing.T) {
	mockedAppInfoRepo := new(mocks.MockedAppInfoRepo)
	var id = "0"
	mockedAppInfoRepo.On("FindOneById", mock.Anything, id).Return(helper.QueryResult{
		nil,
		nil,
	})
	repoContainer := repository.RepositoryContainer{
		AppInfoRepo: mockedAppInfoRepo,
	}
	c := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: NewResolver(&repoContainer)})))

	query := `{ 
		application(id: "0") {
			id
		}}`

	response, _ := c.RawPost(query)

	responseErrors := extractResponseErrors(response)

	require.Nil(t, response.Data)
	require.Equal(t, "application", responseErrors[0]["path"].([]interface{})[0])
	require.Equal(t, "Application with id 0 not found", responseErrors[0]["message"])

}

func extractResponseErrors(response *client.Response) []map[string]interface{} {
	var responseErrors []map[string]interface{}
	err := json.Unmarshal(response.Errors, &responseErrors)
	if err != nil {
		log.Panicf("Error unmarshalling errors into Json structure: %v", err)
	}
	return responseErrors
}
