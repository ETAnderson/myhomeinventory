package inventory

import (
    "fmt"
    "os"
    "time"
)

// EnsureTables verifies required tables exist and are properly structured.
func (d *Database) EnsureTables() {
    tables := []struct {
        Name         string
        CreateStmt   string
        ExpectedCols []string
    }{
        {
            Name: "inventory_item",
            CreateStmt: `
                CREATE TABLE inventory_item (
                    id INT AUTO_INCREMENT PRIMARY KEY,
                    item_name VARCHAR(255) NOT NULL,
                    itemQTY INT NOT NULL,
                    minimumQTY INT NOT NULL,
                    itemUsedToDate INT NOT NULL DEFAULT 0,
                    item_type_id INT,
                    item_substitution_id INT,
                    createDate DATETIME DEFAULT CURRENT_TIMESTAMP,
                    lastModifiedDate DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                    FOREIGN KEY (item_type_id) REFERENCES item_type(id),
                    FOREIGN KEY (item_substitution_id) REFERENCES item_substitution(id)
                );
            `,
            ExpectedCols: []string{"id", "item_name", "itemQTY", "minimumQTY", "itemUsedToDate", "item_type_id", "item_substitution_id", "createDate", "lastModifiedDate"},
        },
        {
            Name: "item_type",
            CreateStmt: `
                CREATE TABLE item_type (
                    id INT AUTO_INCREMENT PRIMARY KEY,
                    type_name VARCHAR(255) NOT NULL UNIQUE
                );
            `,
            ExpectedCols: []string{"id", "type_name"},
        },
        {
            Name: "item_substitution",
            CreateStmt: `
                CREATE TABLE item_substitution (
                    id INT AUTO_INCREMENT PRIMARY KEY,
                    substitution_name VARCHAR(255) NOT NULL UNIQUE
                );
            `,
            ExpectedCols: []string{"id", "substitution_name"},
        },
    }

    for _, table := range tables {
        d.checkAndCreateTable(table.Name, table.CreateStmt, table.ExpectedCols)
    }
}

// checkAndCreateTable verifies a table exists and matches the expected structure.
func (d *Database) checkAndCreateTable(tableName, createStmt string, expectedCols []string) {
    row := d.conn.QueryRow(`
        SELECT COUNT(*)
        FROM information_schema.tables
        WHERE table_schema = DATABASE()
        AND table_name = ?
    `, tableName)

    var count int
    var err error
    err = row.Scan(&count)
    exists := (err == nil && count > 0)

    if !exists {
        fmt.Printf("Table '%s' does not exist.\n", tableName)
        if confirm(fmt.Sprintf("Create table '%s'? (yes/no): ", tableName)) {
            if _, err := d.conn.Exec(createStmt); err != nil {
                fmt.Printf("Failed to create table '%s': %v\n", tableName, err)
                os.Exit(1)
            }
            fmt.Printf("Table '%s' created successfully.\n", tableName)
        } else {
            fmt.Printf("Table '%s' creation aborted by user.\n", tableName)
            os.Exit(1)
        }
        return
    }

    if !d.ValidateTableStructure(tableName, expectedCols) {
        fmt.Printf("Table '%s' exists but is improperly structured.\n", tableName)
        if confirm(fmt.Sprintf("Archive and recreate table '%s'? (yes/no): ", tableName)) {
            timestamp := time.Now().Format("20060102_150405")
            archiveName := fmt.Sprintf("%s_%s", tableName, timestamp)
            renameStmt := fmt.Sprintf("RENAME TABLE %s TO %s", tableName, archiveName)
            if _, err := d.conn.Exec(renameStmt); err != nil {
                fmt.Printf("Failed to archive table '%s': %v\n", tableName, err)
                os.Exit(1)
            }
            fmt.Printf("Table '%s' archived as '%s'.\n", tableName, archiveName)

            if _, err := d.conn.Exec(createStmt); err != nil {
                fmt.Printf("Failed to create new table '%s': %v\n", tableName, err)
                os.Exit(1)
            }
            fmt.Printf("Table '%s' created successfully.\n", tableName)
        } else {
            fmt.Printf("Table '%s' recreation aborted by user.\n", tableName)
            os.Exit(1)
        }
    } else {
        fmt.Printf("Table '%s' exists and structure is valid.\n", tableName)
    }
}
