package config

import (
	"os"
	"strconv"
)

func ternary(value string, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

var (
	// IsMaster master
	IsMaster = os.Getenv("MASTER") == "1"

	// Verbose verbose
	Verbose = ternary(os.Getenv("VERBOSE"), "1") == "1"

	// HTTPPort - http port to run
	HTTPPort, _  = strconv.Atoi(ternary(os.Getenv("HTTP_PORT"), "8088"))
	DBConnection = ternary(os.Getenv("DB_CONNECTION"), "postgresql://postgres:9406715@localhost:5432/personalDB?sslmode=disable")
)
