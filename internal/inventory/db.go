package inventory

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql" // MySQL driver
)

type Database struct {
    conn *sql.DB
}

func NewDatabase() *Database {
    return &Database{}
}

func (d *Database) Boot() {
    if d.IsRunning() {
        fmt.Println("Database is already running.")
        return
    }

    dbUser := GetEnv("DB_USER")
    dbPassword := GetEnv("DB_PASSWORD")
    dbHost := GetEnv("DB_HOST")
    dbPort := GetEnv("DB_PORT")
    dbName := GetEnv("DB_NAME")

    // ✅ Added parseTime=true to properly scan timestamps into time.Time
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

    var err error
    d.conn, err = sql.Open("mysql", dsn)
    if err != nil {
        panic(fmt.Sprintf("Error opening database: %v", err))
    }

    if err := d.conn.Ping(); err != nil {
        panic(fmt.Sprintf("Error pinging database: %v", err))
    }

    // ✅ New: Log current schema/database name
    var currentSchema string
    err = d.conn.QueryRow("SELECT DATABASE()").Scan(&currentSchema)
    if err != nil {
        panic(fmt.Sprintf("Error getting current database: %v", err))
    }

    fmt.Println("✅ Database connected successfully.")
    fmt.Printf("📚 Connected to schema (database): %s\n", currentSchema)
    fmt.Printf("📦 Using table: item_inventory\n") // Hardcoded for now
}

func (d *Database) Shutdown() {
    if d.conn != nil {
        d.conn.Close()
        d.conn = nil
        fmt.Println("Database connection closed.")
    }
}

func (d *Database) IsRunning() bool {
    return d.conn != nil
}

func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
    return d.conn.QueryRow(query, args...)
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
    return d.conn.Query(query, args...)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
    return d.conn.Exec(query, args...)
}
