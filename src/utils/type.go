package utils

func String(s string) *string {
	return &s
}

func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func Int(s int) *int {
	return &s
}

func IntVal(s *int) int {
	if s == nil {
		return 0
	}
	return *s
}


func Bool(s bool) *bool {
	return &s
}