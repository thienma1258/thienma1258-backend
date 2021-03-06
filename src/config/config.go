package config

import (
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
		log.Info("Error getting env, not comming through %v", err)
	} else {
		log.Printf("we getting env variable")
	}
	//Get all env variables
	log.Printf("app running env %v", os.Environ())
}

var (
	// IsMaster master
	IsMaster = os.Getenv("MASTER") == "1"

	// Verbose verbose
	Verbose     = ternary(os.Getenv("VERBOSE"), "1") == "1"
	ServiceName = ternary(os.Getenv("SERVICE_NAME"), "API pham ngoc dong")

	// HTTPPort - http port to run
	HTTPPort, _  = strconv.Atoi(ternary(os.Getenv("HTTP_PORT"), "8088"))
	HTTPSPort, _ = strconv.Atoi(ternary(os.Getenv("HTTPS_PORT"), "8443"))

	DBConnection = ternary(os.Getenv("PERSONAL_DB_CONNECTION"), "postgresql://postgres:9406715@localhost:5432/personalDB?sslmode=disable")
	RedisAddr    = ternary(os.Getenv("REDIS_SERVICE_ADDR"), "127.0.0.1:6379")
	SSLKeyBase64 = os.Getenv("SSL_KEY_64")

	BloggerAPIKey = ternary(os.Getenv("BLOGGER_API_KEY"),"")
	NewApiKey = ternary(os.Getenv("NEW_API_KEY"),"")


)
