package inventory

import (
    "fmt"
)

// InsertItem inserts a new inventory item into the database.
func InsertItem(db *Database, item InventoryItem) (int64, error) {
    query := `
        INSERT INTO inventory_item 
        (item_name, itemQTY, minimumQTY, itemUsedToDate, item_type_id, item_substitution_id)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    result, err := db.Exec(query,
        item.ItemName,
        item.ItemQTY,
        item.MinimumQTY,
        item.ItemUsedToDate,
        item.ItemTypeID,
        item.ItemSubstitutionID,
    )
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}

// GetItemList retrieves a list of inventory items with their type and substitution names.
func GetItemList(db *Database, limit int, itemType string, underMinimum bool) ([]InventoryItemWithDetails, error) {
    query := `
        SELECT 
            i.id, 
            i.item_name, 
            i.itemQTY, 
            i.minimumQTY, 
            i.itemUsedToDate,
            t.type_name,
            s.substitution_name,
            i.createDate, 
            i.lastModifiedDate
        FROM inventory_item i
        LEFT JOIN item_type t ON i.item_type_id = t.id
        LEFT JOIN item_substitution s ON i.item_substitution_id = s.id
        WHERE 1=1
    `
    args := []interface{}{}

    if itemType != "" {
        query += " AND t.type_name = ?"
        args = append(args, itemType)
    }

    if underMinimum {
        query += " AND i.itemQTY < i.minimumQTY"
    }

    query += " ORDER BY i.id ASC"

    if limit > 0 {
        query += " LIMIT ?"
        args = append(args, limit)
    }

    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    items := []InventoryItemWithDetails{}
    for rows.Next() {
        var item InventoryItemWithDetails
        err := rows.Scan(
            &item.ID,
            &item.ItemName,
            &item.ItemQTY,
            &item.MinimumQTY,
            &item.ItemUsedToDate,
            &item.ItemTypeName,
            &item.ItemSubstitutionName,
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

// GetItemTypes retrieves all item types from the database.
func GetItemTypes(db *Database) ([]ItemType, error) {
    query := `SELECT id, type_name FROM item_type ORDER BY type_name ASC`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var types []ItemType
    for rows.Next() {
        var t ItemType
        if err := rows.Scan(&t.ID, &t.Name); err != nil {
            return nil, err
        }
        types = append(types, t)
    }
    return types, nil
}

// GetItemSubstitutions retrieves all item substitutions from the database.
func GetItemSubstitutions(db *Database) ([]ItemSubstitution, error) {
    query := `SELECT id, substitution_name FROM item_substitution ORDER BY substitution_name ASC`
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var substitutions []ItemSubstitution
    for rows.Next() {
        var s ItemSubstitution
        if err := rows.Scan(&s.ID, &s.Name); err != nil {
            return nil, err
        }
        substitutions = append(substitutions, s)
    }
    return substitutions, nil
}

// UpdateItemQty updates the quantity and used count of an inventory item.
func UpdateItemQty(db *Database, itemName string, action string) (map[string]interface{}, error) {
    if action != "+" && action != "-" {
        return nil, fmt.Errorf("invalid action: must be + or -")
    }

    var updateQuery string
    if action == "+" {
        updateQuery = `
            UPDATE inventory_item
            SET itemQTY = itemQTY + 1, lastModifiedDate = NOW()
            WHERE item_name = ?
        `
    } else {
        updateQuery = `
            UPDATE inventory_item
            SET itemQTY = itemQTY - 1, itemUsedToDate = itemUsedToDate + 1, lastModifiedDate = NOW()
            WHERE item_name = ?
        `
    }

    _, err := db.Exec(updateQuery, itemName)
    if err != nil {
        return nil, err
    }

    selectQuery := `
        SELECT id, item_name, itemQTY, itemUsedToDate
        FROM inventory_item
        WHERE item_name = ?
    `
    row := db.QueryRow(selectQuery, itemName)

    var id, qty, used int
    var name string
    if err := row.Scan(&id, &name, &qty, &used); err != nil {
        return nil, err
    }

    result := map[string]interface{}{
        "id":             id,
        "itemName":       name,
        "itemQTY":        qty,
        "itemUsedToDate": used,
    }
    return result, nil
}
