package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCoreUtils_AddMapToMap_NilToEmpty(t *testing.T) {
	var src map[string]any
	var dst = map[string]any{}
	MergeMapToMap(src, dst)
	require.Empty(t, dst)
}

func TestCoreUtils_AddMapToMap_NilToNotEmpty(t *testing.T) {
	var src map[string]any
	var dst = map[string]any{"k": "v"}
	MergeMapToMap(src, dst)

	require.NotNil(t, dst)
	require.Equal(t, "v", dst["k"])
}

func TestCoreUtils_AddMapToMap_NotEmptyToEmpty(t *testing.T) {
	var src = map[string]any{"k": "v"}
	var dst = map[string]any{}
	MergeMapToMap(src, dst)

	require.NotNil(t, dst)
	require.Equal(t, "v", dst["k"])
}

func TestCoreUtils_AddMapToMap_NotEmptyToNotEmpty(t *testing.T) {
	var src = map[string]any{"k": "v"}
	var dst = map[string]any{"e": "f"}
	MergeMapToMap(src, dst)

	require.NotNil(t, dst)
	require.Equal(t, "v", dst["k"])
	require.Equal(t, "f", dst["e"])
}
