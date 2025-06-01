package inventory

import "time"

// InventoryItem represents a record in the inventory_item table.
type InventoryItem struct {
    ID                  int       `json:"id"`
    ItemName            string    `json:"itemName"`
    ItemQTY             int       `json:"itemQTY"`
    MinimumQTY          int       `json:"minimumQTY"`
    ItemUsedToDate      int       `json:"itemUsedToDate"`
    ItemTypeID          int       `json:"itemTypeID"`
    ItemSubstitutionID  int       `json:"itemSubstitutionID"`
    ItemExpirationPeriod int      `json:"itemExpirationPeriod"`
    ItemTotalTossed     int       `json:"itemTotalTossed"`
    CreateDate          time.Time `json:"createDate"`
    LastModifiedDate    time.Time `json:"lastModifiedDate"`
}

// InventoryItemWithDetails represents an inventory item with type and substitution names for display purposes.
type InventoryItemWithDetails struct {
    ID                   int       `json:"id"`
    ItemName             string    `json:"itemName"`
    ItemQTY              int       `json:"itemQTY"`
    MinimumQTY           int       `json:"minimumQTY"`
    ItemUsedToDate       int       `json:"itemUsedToDate"`
    ItemTotalTossed      int       `json:"itemTotalTossed"`
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
