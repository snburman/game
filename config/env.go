package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Environment variables
type Vars struct {
	SERVER_URL    string
	CLIENT_ID     string
	CLIENT_SECRET string
	MONGO_URI     string
}

// Env() returns Vars struct of environment variables
func Env() Vars {
	// Load if not a test. This isn't required during testing.
	if flag.Lookup("test.v") == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading environment variables")
		}
	}

	return Vars{
		SERVER_URL:    os.Getenv("SERVER_URL"),
		CLIENT_ID:     os.Getenv("CLIENT_ID"),
		CLIENT_SECRET: os.Getenv("CLIENT_SECRET"),
		MONGO_URI:     os.Getenv("MONGO_URI"),
	}
}
