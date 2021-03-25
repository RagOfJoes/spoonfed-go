package util

import (
	"os"
	"strconv"

	"github.com/RagOfJoes/spoonfed-go/pkg/logger"
)

// GetAssert will return the env or panic if it is not present
func GetAssert(k string) string {
	v := os.Getenv(k)
	if v == "" {
		logger.Panic("ENV missing, key: " + k)
	}
	return v
}

// GetAssertBool will return the env as boolean or panic if it is not present
func GetAssertBool(k string) bool {
	v := os.Getenv(k)
	if v == "" {
		logger.Panic("ENV missing, key: " + k)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		logger.Panic("ENV err: [" + k + "]" + err.Error())
	}
	return b
}

// GetAssertInt32 will return the env as int32 or panic if it is not present
func GetAssertInt32(k string) int {
	v := os.Getenv(k)
	if v == "" {
		logger.Panic("ENV missing, key: " + k)
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		logger.Panic("ENV err: [" + k + "]" + err.Error())
	}
	return int(i)
}

// GetAssertInt64 will return the env as int64 or panic if it is not present
func GetAssertInt64(k string) int64 {
	v := os.Getenv(k)
	if v == "" {
		logger.Panic("ENV missing, key: " + k)
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logger.Panic("ENV err: [" + k + "]" + err.Error())
	}
	return i
}

// GetEnv is a simple helper function to read an environment or return a default value
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
