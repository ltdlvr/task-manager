package config

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
)

func fatal(prfx string, envVar string, val any) {
	log.Fatalf("ERROR: %s is invalid, %s = %+v", prfx, envVar, val)
}

func isValidPort(port string) bool {
	val, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return val >= 1 && val <= 65535
}

func getDatabaseURL() string {
	user := os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	pswd := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	sslMode := os.Getenv("DB_SSL_MODE")

	if user == "" {
		fatal("Database user", "DB_USER", user)
	}
	if pswd == "" {
		fatal("Database password", "DB_PASSWORD", pswd)
	}
	if host == "" {
		fatal("Database host", "DB_HOST", host)
	}
	if name == "" {
		fatal("Database name", "DB_NAME", name)
	}
	if !isValidPort(port) {
		fatal("Database port", "DB_PORT", port)
	}
	if sslMode != "disable" && sslMode != "enable" {
		fatal("Database SSL mode", "DB_SSL_MODE", sslMode)
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pswd, host, port, name, sslMode,
	)
}
