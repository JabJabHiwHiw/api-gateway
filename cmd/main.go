package main

import (
	"context"
	"food-service/api"
	"log"
	"net/http"
    "os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

func enableCORS(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        h.ServeHTTP(w, r)
    })
}

func main() {
    clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
    if clerkSecretKey == "" {
        log.Fatal("CLERK_SECRET_KEY is not set")
    }
    clerk.SetKey(clerkSecretKey)

    mux := http.NewServeMux()
    mux.HandleFunc("/food/menus", api.GetMenusHandler)
    mux.HandleFunc("/food/menu", api.MenuItemHandler)
    mux.HandleFunc("/food/menu/", api.MenuItemHandler)
    mux.HandleFunc("/food/menu/search", api.SearchMenuHandler)
    mux.HandleFunc("/food/ingredients", api.GetIngredientsHandler)
    mux.HandleFunc("/food/ingredient/", api.IngredientItemHandler)

    mux.Handle("/food/fridge", authMiddleware(http.HandlerFunc(api.GetFridgeHandler)))
    mux.Handle("/food/fridge/item", authMiddleware(http.HandlerFunc(api.FridgeItemHandler)))
    mux.Handle("/food/fridge/item/", authMiddleware(http.HandlerFunc(api.FridgeItemHandler)))
    mux.Handle("/food/fridge/expiring", authMiddleware(http.HandlerFunc(api.GetExpiringFridgeItemsHandler)))

    loggedMux := enableCORS(mux)

    log.Println("Food Service is running on port 8080 with CORS enabled...")
    if err := http.ListenAndServe(":8080", loggedMux); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func authMiddleware(next http.Handler) http.Handler {
    return clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims, ok := clerk.SessionClaimsFromContext(r.Context())
        if !ok {
            http.Error(w, `{"access": "unauthorized"}`, http.StatusUnauthorized)
            return
        }

        usr, err := user.Get(r.Context(), claims.Subject)
        if err != nil {
            http.Error(w, `{"error": "failed to retrieve user"}`, http.StatusInternalServerError)
            return
        }

        ctx := context.WithValue(r.Context(), "user_id", usr.ID)

        next.ServeHTTP(w, r.WithContext(ctx))
    }))
}