package handlers

import (
	"context"
	"encoding/json"
	"github.com/Andision/RoundWorld/framework/web/structures"
	"github.com/Andision/RoundWorld/games"
	"log"
	"net/http"
)

type createTableData struct {
	GameType string `json:"game_type"`
}

func CreateTableHandler(ctx context.Context, lounge structures.Lounge, gamesConfig *games.ConfigGames, w http.ResponseWriter, r *http.Request) {
	var data createTableData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameConfig, exists := gamesConfig.GetAvailableGames()[data.GameType]
	if !exists {
		log.Printf("Game type %s not found", data.GameType)
		http.Error(w, "Game type not found", http.StatusBadRequest)
		return
	}

	executor, ok := ctx.Value("username").(string)
	if !ok {
		log.Printf("Username not found in context")
		http.Error(w, "username not found in ctx", http.StatusBadRequest)
		return
	}

	table := structures.NewTable(gameConfig, executor)
	go table.Run()

	err = lounge.AddTable(table)
	if err != nil {
		log.Printf("Failed to add table: %v", err)
		http.Error(w, "Failed to create table", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// 创建包含 table ID 的响应结构体
	response := struct {
		TableID string `json:"table_id"`
	}{
		TableID: table.GetTableId(),
	}

	// 将响应结构体编码为 JSON 并写入响应
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
