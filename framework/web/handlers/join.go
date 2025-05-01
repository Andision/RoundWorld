package handlers

import (
	"context"
	"github.com/Andision/RoundWorld/framework/web/structures"
	"log"
	"net/http"
)

func JoinHandler(ctx context.Context, lounge structures.Lounge, w http.ResponseWriter, r *http.Request) {
	log.Println("JoinHandler called")
	tableId := r.URL.Query().Get("table_id")

	table := lounge.GetTableById(tableId)
	if table == nil {
		log.Println("TableImpl not found")
		http.Error(w, "TableImpl not found", http.StatusNotFound)
		return
	}

	table.Join(ctx, w, r)
}
