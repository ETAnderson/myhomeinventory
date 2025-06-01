package server

import (
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "strconv"

    "myhomeinventory/internal/inventory"
)

// serveHome serves the home page (not used since we switched to dynamic form rendering).
func serveHome(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "templates/index.html")
}

// makeHandleItems returns an HTTP handler that retrieves the list of inventory items.
func makeHandleItems(db *inventory.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        items, err := inventory.GetItemList(db, 0, "", false)
        if err != nil {
            fmt.Println("Failed to get items:", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(items)
    }
}

// makeHandleUpdateItem returns an HTTP handler that updates the quantity of an inventory item.
func makeHandleUpdateItem(db *inventory.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }
        itemName := r.FormValue("itemName")
        action := r.FormValue("action")

        result, err := inventory.UpdateItemQty(db, itemName, action)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(result)
    }
}

// makeHandleAddItemForm returns an HTTP handler that serves the Add Item form with dynamic dropdowns.
func makeHandleAddItemForm(db *inventory.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        itemTypes, err := inventory.GetItemTypes(db)
        if err != nil {
            fmt.Println("Failed to fetch item types:", err)
            http.Error(w, "Failed to load form", http.StatusInternalServerError)
            return
        }

        itemSubstitutions, err := inventory.GetItemSubstitutions(db)
        if err != nil {
            fmt.Println("Failed to fetch item substitutions:", err)
            http.Error(w, "Failed to load form", http.StatusInternalServerError)
            return
        }

        tmpl, err := template.ParseFiles("templates/index.html") // ✅ Corrected here
        if err != nil {
            fmt.Println("Failed to parse template:", err)
            http.Error(w, "Failed to load form", http.StatusInternalServerError)
            return
        }

        data := struct {
            ItemTypes         []inventory.ItemType
            ItemSubstitutions []inventory.ItemSubstitution
        }{
            ItemTypes:         itemTypes,
            ItemSubstitutions: itemSubstitutions,
        }

        if err := tmpl.Execute(w, data); err != nil {
            fmt.Println("Failed to render template:", err)
            http.Error(w, "Failed to load form", http.StatusInternalServerError)
        }
    }
}

// makeHandleAddItem returns an HTTP handler that adds a new inventory item.
func makeHandleAddItem(db *inventory.Database) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
            return
        }

        itemName := r.FormValue("itemName")
        qtyStr := r.FormValue("itemQTY")
        minQtyStr := r.FormValue("minimumQTY")
        itemTypeIDStr := r.FormValue("itemTypeID")
        itemSubstitutionIDStr := r.FormValue("itemSubstitutionID")

        qty, err := strconv.Atoi(qtyStr)
        if err != nil {
            fmt.Println("Invalid quantity:", qtyStr)
            http.Error(w, "Invalid quantity", http.StatusBadRequest)
            return
        }

        minQty, err := strconv.Atoi(minQtyStr)
        if err != nil {
            fmt.Println("Invalid minimum quantity:", minQtyStr)
            http.Error(w, "Invalid minimum quantity", http.StatusBadRequest)
            return
        }

        itemTypeID, err := strconv.Atoi(itemTypeIDStr)
        if err != nil {
            fmt.Println("Invalid item type ID:", itemTypeIDStr)
            http.Error(w, "Invalid item type selection", http.StatusBadRequest)
            return
        }

        itemSubstitutionID, err := strconv.Atoi(itemSubstitutionIDStr)
        if err != nil {
            fmt.Println("Invalid item substitution ID:", itemSubstitutionIDStr)
            http.Error(w, "Invalid item substitution selection", http.StatusBadRequest)
            return
        }

        newItem := inventory.InventoryItem{
            ItemName:           itemName,
            ItemQTY:            qty,
            MinimumQTY:         minQty,
            ItemUsedToDate:     0,
            ItemTypeID:         itemTypeID,
            ItemSubstitutionID: itemSubstitutionID,
        }

        id, err := inventory.InsertItem(db, newItem)
        if err != nil {
            fmt.Println("Failed to insert item:", err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        resp := map[string]interface{}{
            "id": id,
        }
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    }
}
