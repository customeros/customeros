package utils

import "strings"

func FilterEmpty(vals []string) []string {
	var result []string
	for _, val := range vals {
		if val != "" {
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
