package writescript

// IsValueInList iterate over an array of strings and check if a value is equal
func IsValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
