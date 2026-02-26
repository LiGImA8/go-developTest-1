package config

import (
	"fmt"
	"os"
)

func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func MustDSN() string {
	host := GetEnv("MYSQL_HOST", "127.0.0.1")
	port := GetEnv("MYSQL_PORT", "3306")
	user := GetEnv("MYSQL_USER", "minigate")
	pass := GetEnv("MYSQL_PASSWORD", "minigate")
	name := GetEnv("MYSQL_DB", "minigate")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
}
