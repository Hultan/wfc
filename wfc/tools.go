package wfc

func replaceCharInString(original, replace string, i int) string {
	if i < 0 || i > len(original) {
		return original
	}
	return original[0:i] + replace + original[i+1:]
}

func getKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
