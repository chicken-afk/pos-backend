package utils

import "strings"

func ConvertSliceToString(slice []string) string {
	result := ""
	for i, str := range slice {
		if i > 0 {
			result += ","
		}
		result += str
	}
	return result
}

func ConvertStringToSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
