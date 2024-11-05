package api

import (
	"context"
	"encoding/json"
	"food-service/grpc_clients"
	"food-service/proto"
	"net/http"
)

func GetMenusHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    client, conn, err := grpc_clients.NewMenuServiceClient("localhost:50051")
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

    w.Header().Set("Content-Type", "application/json")
    w.Write(prettyJSON)
}

func MenuItemHandler(w http.ResponseWriter, r *http.Request) {
	client, conn, err := grpc_clients.NewMenuServiceClient("localhost:50051")
	if err != nil {
		http.Error(w, "Failed to connect to Menu Service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	if r.Method == http.MethodGet {
		grpcResponse, err := client.GetMenuItem(context.Background(), &proto.MenuItemRequest{Id: r.URL.Path[len("/menu/"):]})
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
		grpcResponse, err := client.DeleteMenuItem(context.Background(), &proto.MenuItemRequest{Id: r.URL.Path[len("/menu/"):]})
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