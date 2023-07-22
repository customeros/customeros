package utils

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/net/context"
)

// Deprecated, use similar method from common module
func ExtractAllRecordsAsString(ctx context.Context, result neo4j.ResultWithContext, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]string, 0)
	for _, v := range records {
		output = append(output, v.Values[0].(string))
	}
	return output, nil
}
