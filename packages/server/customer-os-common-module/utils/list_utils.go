package utils

import (
	"fmt"
	"strings"
)

func RemoveEmpties(arr []string) []string {
	var result []string
	for _, val := range arr {
		if val != "" {
			result = append(result, val)
		}
	}
	return result
}

func AddToListIfNotExists(arr []string, valueToAdd string) []string {
	if !Contains(arr, valueToAdd) {
		arr = append(arr, valueToAdd)
	}
	return arr
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

// Deprecated: use LowercaseSliceOfStrings instead
func LowercaseStrings(arr []string) {
	for i, s := range arr {
		arr[i] = strings.ToLower(s)
	}
}

func LowercaseSliceOfStrings(arr []string) []string {
	var result []string
	for _, s := range arr {
		result = append(result, strings.ToLower(s))
	}
	return result
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ContainsElement[T comparable](slice []T, value T) bool {
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

func StringSlicesEqualIgnoreOrder(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	// Create maps to count occurrences of each string
	count1 := make(map[string]int)
	count2 := make(map[string]int)

	// Count occurrences in both slices
	for _, s := range slice1 {
		count1[s]++
	}
	for _, s := range slice2 {
		count2[s]++
	}

	// Compare the counts
	for s, c := range count1 {
		if count2[s] != c {
			return false
		}
	}

	for s, c := range count2 {
		if count1[s] != c {
			return false
		}
	}

	return true
}
