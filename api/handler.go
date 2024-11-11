package api

import (
	"context"
	"encoding/json"
	"fmt"
	"food-service/grpc_clients"
	"food-service/proto"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func GetMenusHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    client, conn, err := grpc_clients.NewMenuServiceClient("menu-service:50051")
    if err != nil {
        http.Error(w, "Failed to connect to Menu Service", http.StatusInternalServerError)
        return
    }
    defer conn.Close()

    grpcResponse, err := client.GetMenus(context.Background(), &proto.Empty{})
    if err != nil {
        http.Error(w, "Failed to retrieve menus", http.StatusInternalServerError)
        return
    }

    prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
    if err != nil {
        http.Error(w, "Failed to format response", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
    w.Header().Set("Content-Type", "application/json")
    w.Write(prettyJSON)
}

func MenuItemHandler(w http.ResponseWriter, r *http.Request) {
	client, conn, err := grpc_clients.NewMenuServiceClient("menu-service:50051")
	if err != nil {
		http.Error(w, "Failed to connect to Menu Service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	if r.Method == http.MethodGet {
		grpcResponse, err := client.GetMenuItem(context.Background(), &proto.MenuItemRequest{Id: r.URL.Path[len("/food/menu/"):]})
		if err != nil {
			http.Error(w, "Failed to retrieve menus", http.StatusInternalServerError)
			return
		}
		prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
		if err != nil {
			http.Error(w, "Failed to format response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(prettyJSON)
	}

	if r.Method == http.MethodPost {
		var menuItem proto.MenuItem
		if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		grpcResponse, err := client.CreateMenuItem(context.Background(), &menuItem)
		if err != nil {
			http.Error(w, "Failed to add menu item", http.StatusInternalServerError)
			return
		}
		prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
		if err != nil {
			http.Error(w, "Failed to format response", http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(prettyJSON)
	}

	if r.Method == http.MethodDelete {
		grpcResponse, err := client.DeleteMenuItem(context.Background(), &proto.MenuItemRequest{Id: r.URL.Path[len("/food/menu/"):]})
		if err != nil {
			http.Error(w, "Failed to delete menu item", http.StatusInternalServerError)
			return
		}
		prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
		if err != nil {
			http.Error(w, "Failed to format response", http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(prettyJSON)
	}

	if r.Method == http.MethodPut {
		var menuItem proto.MenuItem
		if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		grpcResponse, err := client.UpdateMenuItem(context.Background(), &menuItem)
		if err != nil {
			http.Error(w, "Failed to update menu item", http.StatusInternalServerError)
			return
		}
		prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
		if err != nil {
			http.Error(w, "Failed to format response", http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(prettyJSON)
	}
}

func SearchMenuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client, conn, err := grpc_clients.NewMenuServiceClient("menu-service:50051")
	if err != nil {
		http.Error(w, "Failed to connect to Menu Service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	grpcResponse, err := client.SearchMenus(context.Background(), &proto.MenuSearchRequest{SearchTerm: r.URL.Query().Get("query")})
	if err != nil {
		http.Error(w, "Failed to search menu", http.StatusInternalServerError)
		return
	}

	prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
	if err != nil {
		http.Error(w, "Failed to format response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(prettyJSON)
}


func GetFridgeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := r.Context().Value("user_id").(string)
    if !ok || userID == "" {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("user-id", userID))

	client, conn, err := grpc_clients.NewFridgeItemServiceClient("fridge-service:50052")
	if err != nil {
		http.Error(w, "Failed to connect to Fridge Service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	fmt.Println(&proto.FridgeRequest{UserId: userID})
	grpcResponse, err := client.GetFridge(ctx, &proto.FridgeRequest{UserId: userID})
	if err != nil {
		http.Error(w, "Failed to retrieve fridge items", http.StatusInternalServerError)
		return
	}

	prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
	if err != nil {
		http.Error(w, "Failed to format response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(prettyJSON)
}

func FridgeItemHandler(w http.ResponseWriter, r *http.Request) {
	client, conn, err := grpc_clients.NewFridgeItemServiceClient("fridge-service:50052")
	if err != nil {
		http.Error(w, "Failed to connect to Fridge Service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	userID, ok := r.Context().Value("user_id").(string)
    if !ok || userID == "" {
        http.Error(w, "User ID not found in context", http.StatusUnauthorized)
        return
    }

	if r.Method == http.MethodGet {
		grpcResponse, err := client.GetFridgeItem(context.Background(), &proto.FridgeItemRequest{Id: r.URL.Path[len("/food/fridge/item/"):],})
		if err != nil {
			http.Error(w, "Failed to retrieve fridge item", http.StatusInternalServerError)
			return
		}
		prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
		if err != nil {
			http.Error(w, "Failed to format response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(prettyJSON)
	}

	if r.Method == http.MethodPost {
		var fridgeItem proto.FridgeItem

		if err := json.NewDecoder(r.Body).Decode(&fridgeItem); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		fridgeItem.UserId = userID

		grpcResponse, err := client.AddItem(context.Background(), &fridgeItem)
		if err != nil {
			http.Error(w, "Failed to add fridge item", http.StatusInternalServerError)
			return
		}
		prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
		if err != nil {
			http.Error(w, "Failed to format response", http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(prettyJSON)
	}

	if r.Method == http.MethodDelete {
		grpcResponse, err := client.RemoveItem(context.Background(), &proto.FridgeItemRequest{Id: r.URL.Path[len("/food/fridge/item/"):],})
		if err != nil {
			http.Error(w, "Failed to delete fridge item", http.StatusInternalServerError)
			return
		}
		prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
		if err != nil {
			http.Error(w, "Failed to format response", http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(prettyJSON)
	}
}

func GetExpiringFridgeItemsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client, conn, err := grpc_clients.NewFridgeItemServiceClient("fridge-service:50052")
	if err != nil {
		http.Error(w, "Failed to connect to Fridge Service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	grpcResponse, err := client.GetExpiredItems(context.Background(), &proto.FridgeRequest{})
	if err != nil {
		http.Error(w, "Failed to retrieve expiring fridge items", http.StatusInternalServerError)
		return
	}

	prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
	if err != nil {
		http.Error(w, "Failed to format response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(prettyJSON)
}

func GetIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client, conn, err := grpc_clients.NewIngredientServiceClient("fridge-service:50052")
	if err != nil {
		http.Error(w, "Failed to connect to Fridge Service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	grpcResponse, err := client.GetIngredients(context.Background(), &proto.Empty{})
	if err != nil {
		http.Error(w, "Failed to retrieve ingredients", http.StatusInternalServerError)
		return
	}

	prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
	if err != nil {
		http.Error(w, "Failed to format response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(prettyJSON)
}

func IngredientItemHandler(w http.ResponseWriter, r *http.Request) {
	client, conn, err := grpc_clients.NewIngredientServiceClient("fridge-service:50052")
	if err != nil {
		http.Error(w, "Failed to connect to Fridge Service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	if r.Method == http.MethodGet {
		grpcResponse, err := client.GetIngredientItem(context.Background(), &proto.IngredientItemRequest{Id: r.URL.Path[len("/food/ingredient/"):],})
		if err != nil {
			http.Error(w, "Failed to retrieve ingredient", http.StatusInternalServerError)
			return
		}
		prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
		if err != nil {
			http.Error(w, "Failed to format response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(prettyJSON)
	}

	if r.Method == http.MethodPost {
		var ingredient proto.IngredientItem
		if err := json.NewDecoder(r.Body).Decode(&ingredient); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		grpcResponse, err := client.AddIngredient(context.Background(), &ingredient)
		if err != nil {
			http.Error(w, "Failed to add ingredient", http.StatusInternalServerError)
			return
		}
		prettyJSON, err := json.MarshalIndent(grpcResponse, "", "  ")
		if err != nil {
			http.Error(w, "Failed to format response", http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.Write(prettyJSON)
	}
}


