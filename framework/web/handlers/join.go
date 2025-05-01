package handlers

import (
	"context"
	"github.com/Andision/RoundWorld/framework/web/structures"
	"github.com/golang/glog"
	"net/http"
)

func JoinHandler(ctx context.Context, lounge structures.Lounge, w http.ResponseWriter, r *http.Request) {
	glog.Info("JoinHandler called")
	tableId := r.URL.Query().Get("table_id")

	table := lounge.GetTableById(tableId)
	if table == nil {
		glog.ErrorContext(ctx, "TableImpl not found")
		http.Error(w, "TableImpl not found", http.StatusNotFound)
		return
	}

	table.Join(ctx, w, r)
}
