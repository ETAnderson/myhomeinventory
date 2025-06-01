package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "myhomeinventory/internal/inventory"
    "myhomeinventory/server"
)

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
    fmt.Println("Connected to MySQL successfully!")

    err := inventory.EnsureItemInventoryTable(db)
    if err != nil {
        log.Fatalf("Failed to ensure table exists: %v", err)
    }

    fmt.Println("Database is ready.")

    router := server.NewRouter(db)

    // ✨ Dynamically get host and port
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
