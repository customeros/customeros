package utils

func StringFirstNonEmpty(values ...string) string {
	for _, s := range values {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}
