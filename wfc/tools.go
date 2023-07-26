package wfc

import "fmt"

func getKeyPart(key string, part int) string {
	if len(key)%4 != 0 {
		panic(fmt.Sprintf("invalid key: %s", key))
	}
	l := len(key) / 4
	return key[l*part : l*part+l]
}

func replaceCharInString(original, replace string, i int) string {
	replace = Reverse(replace)
	if i < 0 || i > len(original) {
		return original
	}
	k := keySize / 4
	return original[0:i*k] + replace + original[i*k+k:]
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func getKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
