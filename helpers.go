package amo

// Separates strings values with comma
func MultiValuesString(s ...string) string {
	res := ""
	prefix := ""
	for _, si := range s {
		res += prefix + si
		prefix = ","
	}
	return res
}