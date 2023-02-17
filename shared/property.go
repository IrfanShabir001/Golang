package shared

import (
	"os"
	"strconv"
)

/* Configurations  */
type MapPropertySource struct{ Data map[string]interface{} }

func (mps *MapPropertySource) Get(key string) interface{} { return mps.Data[key] }
func (mps *MapPropertySource) GetString(key string) string {
	if os.Getenv(key) != "" {
		return os.Getenv(key)
	}
	return mps.Get(key).(string)
}
func (mps *MapPropertySource) GetInt(key string) int {
	if os.Getenv(key) != "" {
		var value = os.Getenv(key)
		v, err := strconv.Atoi(value)
		if err != nil {
			return mps.Get(key).(int)
		}
		return v
	}

	return mps.Get(key).(int)
}
func (mps *MapPropertySource) GetBool(key string) bool {
	if os.Getenv(key) != "" {
		truth := os.Getenv(key)
		if truth == "true" {
			return true
		}
		return false
	}
	return mps.Get(key).(bool)
}
