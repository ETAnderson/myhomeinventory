package server

import (
    "encoding/json"
    "net/http"
    "strconv"
    "fmt"

    "myhomeinventory/internal/inventory"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "templates/index.html")
}

func makeHandleItems(db *inventory.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        items, err := inventory.GetItemList(db, 0, "", false)
        if err != nil {
            fmt.Println("❌ Failed to get items:", err) 
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(items)
    }
}


func makeHandleUpdateItem(db *inventory.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        itemName := r.FormValue("itemName")
        action := r.FormValue("action") // ✅ Correct: action

        result, err := inventory.UpdateItemQty(db, itemName, action) // ✅ Use action here
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    }
}


func makeHandleAddItem(db *inventory.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("🔥 /item/add called!") 

        if r.Method != http.MethodPost {
            fmt.Println("⛔ Method not allowed:", r.Method) 
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        itemName := r.FormValue("itemName")
        itemType := r.FormValue("itemType")
        qtyStr := r.FormValue("itemQTY")
        minQtyStr := r.FormValue("minimumQTY")

        fmt.Printf("📦 Received form values: itemName=%s, itemType=%s, itemQTY=%s, minimumQTY=%s\n", itemName, itemType, qtyStr, minQtyStr) // ADD THIS

        qty, err := strconv.Atoi(qtyStr)
        if err != nil {
            fmt.Println("❌ Invalid quantity:", qtyStr) 
            http.Error(w, "Invalid quantity", http.StatusBadRequest)
            return
        }

        minQty, err := strconv.Atoi(minQtyStr)
        if err != nil {
            fmt.Println("❌ Invalid minimum quantity:", minQtyStr) 
            http.Error(w, "Invalid minimum quantity", http.StatusBadRequest)
            return
        }

        newItem := inventory.InventoryItem{
            ItemName:       itemName,
            ItemQTY:        qty,
            MinimumQTY:     minQty,
            ItemUsedToDate: 0,
            ItemType:       itemType,
        }

        id, err := inventory.InsertItem(db, newItem)
        if err != nil {
            fmt.Println("❌ Failed to insert item:", err) 
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        fmt.Printf("✅ Successfully inserted item with ID: %d\n", id) 

        resp := map[string]interface{}{
            "id": id,
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }
}

