package env

import (
	"os"
	"strconv"
	"time"
)

// GetString retrieves the value of the environment variable named by the key.
// If the variable is empty or not present, it returns the defaultValue.
func GetString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetInt retrieves the value of the environment variable named by the key and converts it to int64.
// If the variable is empty, not present, or cannot be converted, it returns the defaultValue.
func GetInt(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetBool retrieves the value of the environment variable named by the key and converts it to bool.
// If the variable is empty, not present, or cannot be converted, it returns the defaultValue.
func GetBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

// GetTimeDuration retrieves the value of the environment variable named by the key and converts it to time.Duration.
// If the variable is empty, not present, or cannot be converted, it returns the defaultValue.
func GetTimeDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	durationValue, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return durationValue
}
