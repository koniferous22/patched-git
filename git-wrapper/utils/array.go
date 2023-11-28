package utils

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func ContainsAnyOf(s []string, t []string) bool {
	for _, item := range s {
		if Contains(t, item) {
			return true
		}
	}
	return false
}
