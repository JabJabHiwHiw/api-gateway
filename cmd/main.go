package main

import (
    "log"
    "net/http"
    "food-service/api"
)

func main() {
    http.HandleFunc("/menus", api.GetMenusHandler)
    http.HandleFunc("/menu", api.MenuItemHandler)
    http.HandleFunc("/menu/", api.MenuItemHandler)

    log.Println("Food Service is running on port 8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}