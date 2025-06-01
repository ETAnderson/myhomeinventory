package inventory

import (
    "fmt"
)

func UpdateItemQty(db *Database, itemName string, sign string) (map[string]interface{}, error) {
    if sign != "+" && sign != "-" {
        return nil, fmt.Errorf("invalid sign: must be + or -")
    }

    var updateQuery string
    if sign == "+" {
        updateQuery = `
            UPDATE item_inventory
            SET itemQTY = itemQTY + 1
            WHERE itemName = ?
        `
    } else {
        updateQuery = `
            UPDATE item_inventory
            SET itemQTY = itemQTY - 1, itemUsedToDate = itemUsedToDate + 1
            WHERE itemName = ?
        `
    }

    _, err := db.Exec(updateQuery, itemName)
    if err != nil {
        return nil, err
    }

    selectQuery := `
        SELECT id, itemName, itemQTY, itemUsedToDate
        FROM item_inventory
        WHERE itemName = ?
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
