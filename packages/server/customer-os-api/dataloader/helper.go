package dataloader

import "github.com/graph-gophers/dataloader"

func sortKeys(keys dataloader.Keys) ([]string, map[string]int) {
	var ids []string
	// create a map for remembering the order of keys passed in
	keyOrder := make(map[string]int, len(keys))
	for ix, key := range keys {
		ids = append(ids, key.String())
		keyOrder[key.String()] = ix
	}
	return ids, keyOrder
}
