package inventory

import "time"

type InventoryItem struct {
    ID               int       `json:"id"`
    ItemName         string    `json:"itemName"`
    ItemQTY          int       `json:"itemQTY"`
    MinimumQTY       int       `json:"minimumQTY"`
    ItemUsedToDate   int       `json:"itemUsedToDate"`
    ItemType         string    `json:"itemType"`
    CreateDate       time.Time `json:"createDate"`
    LastModifiedDate time.Time `json:"lastModifiedDate"`
}

func InsertItem(db *Database, item InventoryItem) (int64, error) {
    query := `
        INSERT INTO item_inventory (itemName, itemQTY, minimumQTY, itemUsedToDate, itemType)
        VALUES (?, ?, ?, ?, ?)
    `
    result, err := db.Exec(query, item.ItemName, item.ItemQTY, item.MinimumQTY, item.ItemUsedToDate, item.ItemType)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}

func GetItemList(db *Database, limit int, itemType string, underMinimum bool) ([]InventoryItem, error) {
    query := `
        SELECT id, itemName, itemQTY, minimumQTY, itemUsedToDate, itemType, createDate, lastModifiedDate
        FROM item_inventory
        WHERE 1=1
    `
    args := []interface{}{}

    if itemType != "" {
        query += " AND itemType = ?"
        args = append(args, itemType)
    }

    if underMinimum {
        query += " AND itemQTY < minimumQTY"
    }

    query += " ORDER BY id ASC"

    if limit > 0 {
        query += " LIMIT ?"
        args = append(args, limit)
    }

    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    items := []InventoryItem{}
    for rows.Next() {
        var item InventoryItem
        err := rows.Scan(
            &item.ID,
            &item.ItemName,
            &item.ItemQTY,
            &item.MinimumQTY,
            &item.ItemUsedToDate,
            &item.ItemType,
            &item.CreateDate,
            &item.LastModifiedDate,
        )
        if err != nil {
            return nil, err
        }
        items = append(items, item)
    }

    return items, nil
}
