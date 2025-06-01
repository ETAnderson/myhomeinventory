package inventory

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

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

func GetEnv(key string) string {
    return os.Getenv(key)
}
