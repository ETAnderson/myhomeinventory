package server

import (
    "net/http"
    "myhomeinventory/internal/inventory"
)

// NewRouter creates a new HTTP router with all the application's routes configured.
// It serves static files, API endpoints, and the main application page.
func NewRouter(db *inventory.Database) http.Handler {
    mux := http.NewServeMux()

    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    mux.HandleFunc("/items", makeHandleItems(db))
    mux.HandleFunc("/item/add", makeHandleAddItem(db))
    mux.HandleFunc("/item/update", makeHandleUpdateItem(db))
    mux.HandleFunc("/item/dispose", makeHandleDisposeItem(db)) // <-- New dispose route
    mux.HandleFunc("/", makeHandleAddItemForm(db)) 

    return mux
}
