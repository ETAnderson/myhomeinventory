package inventory

import (
    "fmt"
)

// InsertItem inserts a new inventory item into the database.
func InsertItem(db *Database, item InventoryItem) (int64, error) {
    query := `
        INSERT INTO inventory_item 
        (item_name, itemQTY, minimumQTY, itemUsedToDate, item_type_id, item_substitution_id, item_expiration_period, item_total_tossed)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
    result, err := db.conn.Exec(query,
        item.ItemName,
        item.ItemQTY,
        item.MinimumQTY,
        item.ItemUsedToDate,
        item.ItemTypeID,
        item.ItemSubstitutionID,
        item.ItemExpirationPeriod,
        0,
    )
    if err != nil {
        return 0, err
    }

    itemID, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    if err := insertItemExpirationXref(db, itemID, item.ItemExpirationPeriod, item.ItemQTY); err != nil {
        return 0, err
    }

    return itemID, nil
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
            i.item_total_tossed, -- ✅ Added field here
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

    rows, err := db.conn.Query(query, args...)
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
            &item.ItemTotalTossed, // ✅ Added scan target
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
    rows, err := db.conn.Query(query)
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
    rows, err := db.conn.Query(query)
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

    var itemID, expirationPeriod int
    err := db.conn.QueryRow(`
        SELECT id, item_expiration_period
        FROM inventory_item
        WHERE item_name = ?
    `, itemName).Scan(&itemID, &expirationPeriod)
    if err != nil {
        return nil, err
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

    _, err = db.conn.Exec(updateQuery, itemName)
    if err != nil {
        return nil, err
    }

    if action == "+" {
        if err := insertItemExpirationXref(db, int64(itemID), expirationPeriod, 1); err != nil {
            return nil, err
        }
    } else {
        if err := removeItemExpirationXref(db, int64(itemID), 1); err != nil {
            return nil, err
        }
    }

    row := db.conn.QueryRow(`
        SELECT id, item_name, itemQTY, itemUsedToDate
        FROM inventory_item
        WHERE item_name = ?
    `, itemName)

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

// insertItemExpirationXref inserts expiration tracking rows for a new inventory item.
func insertItemExpirationXref(db *Database, itemID int64, expirationPeriod int, quantity int) error {
    query := `
        INSERT INTO item_expiration_xref (item_id, item_creation_date, item_expiration_date)
        VALUES (?, NOW(), DATE_ADD(NOW(), INTERVAL ? DAY))
    `
    stmt, err := db.conn.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    for i := 0; i < quantity; i++ {
        if _, err := stmt.Exec(itemID, expirationPeriod); err != nil {
            return err
        }
    }

    return nil
}

// removeItemExpirationXref removes the oldest expiration tracking rows for an inventory item.
func removeItemExpirationXref(db *Database, itemID int64, quantity int) error {
    query := `
        DELETE FROM item_expiration_xref
        WHERE id IN (
            SELECT id FROM (
                SELECT id FROM item_expiration_xref
                WHERE item_id = ?
                ORDER BY item_expiration_date ASC
                LIMIT ?
            ) AS sub
        )
    `
    _, err := db.conn.Exec(query, itemID, quantity)
    return err
}

// DisposeItem removes the oldest expiration entry and increments total tossed.
func DisposeItem(db *Database, itemName string) (map[string]interface{}, error) {
    var itemID int
    err := db.conn.QueryRow(`
        SELECT id
        FROM inventory_item
        WHERE item_name = ?
    `, itemName).Scan(&itemID)
    if err != nil {
        return nil, err
    }

    tx, err := db.conn.Begin()
    if err != nil {
        return nil, err
    }

    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)
        } else if err != nil {
            tx.Rollback()
        } else {
            err = tx.Commit()
        }
    }()

    _, err = tx.Exec(`
        DELETE FROM item_expiration_xref
        WHERE id IN (
            SELECT id FROM (
                SELECT id FROM item_expiration_xref
                WHERE item_id = ?
                ORDER BY item_expiration_date ASC
                LIMIT 1
            ) AS sub
        )
    `, itemID)
    if err != nil {
        return nil, err
    }

    _, err = tx.Exec(`
        UPDATE inventory_item
        SET item_total_tossed = item_total_tossed + 1, lastModifiedDate = NOW()
        WHERE id = ?
    `, itemID)
    if err != nil {
        return nil, err
    }

    var tossed, qty int
    err = tx.QueryRow(`
        SELECT item_total_tossed, itemQTY
        FROM inventory_item
        WHERE id = ?
    `, itemID).Scan(&tossed, &qty)
    if err != nil {
        return nil, err
    }

    result := map[string]interface{}{
        "itemTotalTossed": tossed,
        "itemQTY":         qty,
    }

    return result, nil
}
