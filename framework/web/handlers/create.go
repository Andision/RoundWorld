package handlers

import (
	"context"
	"encoding/json"
	"github.com/Andision/RoundWorld/framework/web/structures"
	"github.com/Andision/RoundWorld/games"
	"github.com/golang/glog"
	"net/http"
)

func CreateHandler(ctx context.Context, lounge structures.Lounge, gamesConfig *games.ConfigGames, w http.ResponseWriter, r *http.Request) {
	gameType := r.URL.Query().Get("game_type")

	gameConfig, exists := gamesConfig.GetAvailableGames()[gameType]
	if !exists {
		glog.ErrorContextf(ctx, "Game type %s not found", gameType)
		http.Error(w, "Game type not found", http.StatusBadRequest)
		return
	}

	executor, ok := ctx.Value("username").(string)
	if !ok {
		glog.ErrorContextf(ctx, "Username not found in context")
		http.Error(w, "username not found in ctx", http.StatusBadRequest)
		return
	}

	table := structures.NewTable(gameConfig, executor)
	go table.Run()

	err := lounge.AddTable(table)
	if err != nil {
		glog.ErrorContextf(ctx, "Failed to add table: %v", err)
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
		glog.ErrorContextf(ctx, "Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
