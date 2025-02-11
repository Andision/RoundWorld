package web

import (
	"context"
	"log"
	"net/http"
)

func JoinHandler(ctx context.Context, lounge *Lounge, w http.ResponseWriter, r *http.Request) {
	log.Println("JoinHandler called")
	tableId := r.URL.Query().Get("table_id")

	log.Println("tableId:", tableId)

	table := lounge.GetTable(tableId)
	if table == nil {
		http.Error(w, "Table not found", http.StatusNotFound)
		return
	}

	table.Join(ctx, w, r)
}
