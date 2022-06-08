package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var IsK8 string

// RunMode refers to run mode.
var RunMode string

// ServerPort refers to server port.
var ServerPort string

// PublicKey refers to publickey of EventStoreToken.
var PublicKey string

// EnableAuthentication refers if service to service authentication is enabled.
var EnableAuthentication bool

// EnableOpenTracing set true if opentracing is needed.
var EnableOpenTracing bool

// InitEnvironmentVariables initializes environment variables
func InitEnvironmentVariables() {

	IsK8 = "False"

	RunMode = os.Getenv("RUN_MODE")
	if RunMode == "" {
		RunMode = DEVELOP
	}

	log.Println("RUN MODE:", RunMode)

	if RunMode != PRODUCTION {
		//Load .env file
		err := godotenv.Load()
		if err != nil {
			log.Println("ERROR:", err.Error())
			return
		}
	}
	PublicKey = os.Getenv("PUBLIC_KEY")
	if os.Getenv("ENABLE_AUTHENTICATION") == "" {
		EnableAuthentication = false
	} else {
		if strings.ToLower(os.Getenv("ENABLE_AUTHENTICATION")) == "true" {
			EnableAuthentication = true
		} else {
			EnableAuthentication = false
		}
	}

	if os.Getenv("ENABLE_OPENTRACING") == "" {
		EnableOpenTracing = false
	} else {
		if strings.ToLower(os.Getenv("ENABLE_OPENTRACING")) == "true" {
			EnableOpenTracing = true
		} else {
			EnableOpenTracing = false
		}
	}
	ServerPort = os.Getenv("SERVER_PORT")
	IsK8 = os.Getenv("IS_K8")
}
