package utils

import (
	"fmt"
	"strings"
)

func FilterEmpty(vals []string) []string {
	var result []string
	for _, val := range vals {
		if val != "" {
			result = append(result, val)
		}
	}
	return result
}

func RemoveFromList(arr []string, valueToRemove string) []string {
	var result []string
	for _, val := range arr {
		if val != valueToRemove {
			result = append(result, val)
		}
	}
	return result
}

func RemoveDuplicates(arr []string) []string {
	var result []string
	seen := map[string]bool{}
	for _, val := range arr {
		if !seen[val] {
			result = append(result, val)
			seen[val] = true
		}
	}
	return result
}

func LowercaseStrings(arr []string) {
	for i, s := range arr {
		arr[i] = strings.ToLower(s)
	}
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ContainsAll(sourceSlice, itemsToCheck []string) bool {
	for _, item := range itemsToCheck {
		found := false
		for _, sourceItem := range sourceSlice {
			if sourceItem == item {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func AnySliceToStringSlice(input []any) ([]string, error) {
	result := []string{}
	for _, item := range input {
		str, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("could not convert item to string")
		}
		result = append(result, str)
	}
	return result, nil
}
