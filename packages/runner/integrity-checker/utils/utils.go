package utils

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/pkg/errors"
)

// Deprecated, use similar method from common module
func ExtractSingleRecordFirstValueAsType[T any](ctx context.Context, result neo4j.ResultWithContext, err error) (T, error) {
	value, err := utils.ExtractSingleRecordFirstValue(ctx, result, err)
	if err != nil {
		return *new(T), err
	}

	converted, ok := value.(T)
	if !ok {
		return *new(T), errors.New("invalid type")
	}

	return converted, nil
}
