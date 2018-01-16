package utils

func MapNilCheck(m map[interface{}]interface{}, key interface{}) bool {
	if _, ok := m[key]; ok {
		return false
	}
	return true
}

func IsValidUrl(url string) bool {
	return true
}
