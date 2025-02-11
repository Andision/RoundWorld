package web

import (
	"context"
	"encoding/json"
	"github.com/Andision/RoundWorld/games"
	"net/http"
)

func CreateHandler(ctx context.Context, lounge *Lounge, gamesConfig *games.ConfigGames, w http.ResponseWriter, r *http.Request) {
	gameType := r.URL.Query().Get("game_type")

	gameConfig, exists := gamesConfig.GetAvailableGames()[gameType]
	if !exists {
		http.Error(w, "Game type not found", http.StatusBadRequest)
		return
	}

	executor, ok := ctx.Value("username").(string)
	if !ok {
		http.Error(w, "username not found in ctx", http.StatusBadRequest)
		return
	}

	table := NewTable(gameConfig, executor)
	go table.Run()

	lounge.AddTable(table.GetTableId(), table)

	w.Header().Set("Content-Type", "application/json")
	// 创建包含 table ID 的响应结构体
	response := struct {
		TableID string `json:"table_id"`
	}{
		TableID: table.GetTableId(),
	}

	// 将响应结构体编码为 JSON 并写入响应
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
