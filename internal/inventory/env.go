package inventory

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// LoadEnv loads environment variables from a .env file
// and ensures all required keys are present.
func LoadEnv(requiredKeys []string) {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    missingKeys := []string{}
    for _, key := range requiredKeys {
        if os.Getenv(key) == "" {
            missingKeys = append(missingKeys, key)
        }
    }

    if len(missingKeys) > 0 {
        log.Fatalf("Missing required environment variables: %v", missingKeys)
    }
}

// GetEnv retrieves an environment variable value or panics if not set.
func GetEnv(key string) string {
    value := os.Getenv(key)
    if value == "" {
        log.Fatalf("Environment variable '%s' is required but not set.", key)
    }
    return value
}
