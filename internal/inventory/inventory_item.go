package inventory

import "time"

// InventoryItem represents an item in the inventory.
type InventoryItem struct {
    ID                 int       `json:"id"`
    ItemName           string    `json:"itemName"`
    ItemQTY            int       `json:"itemQTY"`
    MinimumQTY         int       `json:"minimumQTY"`
    ItemUsedToDate     int       `json:"itemUsedToDate"`
    ItemTypeID         int       `json:"itemTypeID"`
    ItemSubstitutionID int       `json:"itemSubstitutionID"`
    CreateDate         time.Time `json:"createDate"`
    LastModifiedDate   time.Time `json:"lastModifiedDate"`
}

// InventoryItemWithDetails represents an item with type and substitution names included.
type InventoryItemWithDetails struct {
    ID                   int       `json:"id"`
    ItemName             string    `json:"itemName"`
    ItemQTY              int       `json:"itemQTY"`
    MinimumQTY           int       `json:"minimumQTY"`
    ItemUsedToDate       int       `json:"itemUsedToDate"`
    ItemTypeName         string    `json:"itemTypeName"`
    ItemSubstitutionName string    `json:"itemSubstitutionName"`
    CreateDate           time.Time `json:"createDate"`
    LastModifiedDate     time.Time `json:"lastModifiedDate"`
}

// ItemType represents a record in the item_type table.
type ItemType struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// ItemSubstitution represents a record in the item_substitution table.
type ItemSubstitution struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
