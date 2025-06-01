package server

import (
    "net/http"
    "myhomeinventory/internal/inventory"
)

func NewRouter(db *inventory.Database) http.Handler {
    mux := http.NewServeMux()

    // Serve static files first
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

    // API routes FIRST
    mux.HandleFunc("/items", makeHandleItems(db))
    mux.HandleFunc("/item/add", makeHandleAddItem(db))
    mux.HandleFunc("/item/update", makeHandleUpdateItem(db))

    // Finally: Serve the SPA
    mux.HandleFunc("/", serveHome)

    return mux
}
