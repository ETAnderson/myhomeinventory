package inventory

import (
    "database/sql"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
)

// Database wraps the sql.DB connection.
type Database struct {
    conn *sql.DB
}

// NewDatabase creates a new instance of Database.
func NewDatabase() *Database {
    return &Database{}
}

// Boot initializes the database connection using environment variables.
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

    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

    var err error
    d.conn, err = sql.Open("mysql", dsn)
    if err != nil {
        panic(fmt.Sprintf("Error opening database: %v", err))
    }

    if err := d.conn.Ping(); err != nil {
        panic(fmt.Sprintf("Error pinging database: %v", err))
    }

    var currentSchema string
    err = d.conn.QueryRow("SELECT DATABASE()").Scan(&currentSchema)
    if err != nil {
        panic(fmt.Sprintf("Error getting current database: %v", err))
    }

    fmt.Println("Database connected successfully.")
    fmt.Printf("Connected to schema: %s\n", currentSchema)
}

// Shutdown cleanly closes the database connection.
func (d *Database) Shutdown() {
    if d.conn != nil {
        d.conn.Close()
        d.conn = nil
        fmt.Println("Database connection closed.")
    }
}

// IsRunning checks if the database connection is active.
func (d *Database) IsRunning() bool {
    return d.conn != nil
}

// QueryRow executes a query that is expected to return at most one row.
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
    return d.conn.QueryRow(query, args...)
}

// Query executes a query that returns rows.
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
    return d.conn.Query(query, args...)
}

// Exec executes a query without returning any rows.
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
    return d.conn.Exec(query, args...)
}

// ValidateTableStructure checks if the table columns match the expected structure.
func (d *Database) ValidateTableStructure(tableName string, expectedCols []string) bool {
    rows, err := d.conn.Query(fmt.Sprintf("DESCRIBE %s", tableName))
    if err != nil {
        fmt.Printf("Failed to describe table '%s': %v\n", tableName, err)
        return false
    }
    defer rows.Close()

    actualCols := []string{}
    for rows.Next() {
        var field, colType, null, key string
        var defaultVal sql.NullString
        var extra string
        if err := rows.Scan(&field, &colType, &null, &key, &defaultVal, &extra); err != nil {
            fmt.Printf("Failed to scan table description: %v\n", err)
            return false
        }
        actualCols = append(actualCols, field)
    }

    if len(actualCols) != len(expectedCols) {
        return false
    }

    for i := range actualCols {
        if actualCols[i] != expectedCols[i] {
            return false
        }
    }

    return true
}
