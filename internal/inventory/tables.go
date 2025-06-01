package inventory

import (
    "fmt"
)

func EnsureItemInventoryTable(db *Database) error {
    query := `
        CREATE TABLE IF NOT EXISTS item_inventory (
            id INT AUTO_INCREMENT PRIMARY KEY,
            itemName VARCHAR(255) NOT NULL,
            itemQTY INT NOT NULL DEFAULT 0,
            minimumQTY INT NOT NULL DEFAULT 0,
            itemUsedToDate INT NOT NULL DEFAULT 0,
            itemType VARCHAR(100) NOT NULL DEFAULT 'general',
            createDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            lastModifiedDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
        ) ENGINE=InnoDB;
    `

    _, err := db.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to ensure table: %w", err)
    }

    return nil
}
