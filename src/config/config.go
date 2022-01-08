package config

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func ternary(value string, defaultValue string) string {
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		log.Printf("we getting env variable")
	}
	//Get all env variables
	fmt.Printf("running with env %v", os.Environ())

}

var (
	// IsMaster master
	IsMaster = os.Getenv("MASTER") == "1"

	// Verbose verbose
	Verbose = ternary(os.Getenv("VERBOSE"), "1") == "1"

	// HTTPPort - http port to run
	HTTPPort, _  = strconv.Atoi(ternary(os.Getenv("HTTP_PORT"), "8088"))
	DBConnection = ternary(os.Getenv("PERSONAL_DB_CONNECTION"), "postgresql://postgres:9406715@localhost:5432/personalDB?sslmode=disable")
)
