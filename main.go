package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "myhomeinventory/internal/inventory"
    "myhomeinventory/server"
)

// main is the entry point of the application.
// It loads environment variables, initializes the database connection,
// ensures required tables exist, sets up the router, and starts the HTTP server.
func main() {
    requiredEnv := []string{
        "DB_USER",
        "DB_PASSWORD",
        "DB_HOST",
        "DB_PORT",
        "DB_NAME",
    }
    inventory.LoadEnv(requiredEnv)

    db := inventory.NewDatabase()
    db.Boot()
    defer db.Shutdown()

    fmt.Println("Database connected successfully.")
    fmt.Println("Connected to MySQL successfully.")

    db.EnsureTables()

    fmt.Println("Database is ready.")

    router := server.NewRouter(db)

    host := os.Getenv("APP_HOST")
    port := os.Getenv("APP_PORT")

    if host == "" {
        host = "localhost"
    }
    if port == "" {
        port = "8080"
    }

    address := fmt.Sprintf("%s:%s", host, port)

    fmt.Println("Starting server on", address)
    fmt.Printf("Server running at: http://%s\n", address)

    log.Fatal(http.ListenAndServe(address, router))
}
