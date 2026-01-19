package utils

// DeepCopyMap 深拷贝map
func DeepCopyMap(original map[string]interface{}) map[string]interface{} {
	copiedMap := make(map[string]interface{})
	for key, value := range original {
		copiedMap[key] = value
	}
	return copiedMap
}
