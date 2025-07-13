package utils

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
