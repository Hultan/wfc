package wfc

func replacePartOfString(s, new string, i int) string {
	if i < 0 || i > len(s) {
		return s
	}
	return s[0:i] + new + s[i+1:]
}

func getKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
