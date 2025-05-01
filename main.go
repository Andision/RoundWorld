package main

import (
	"context"
	"github.com/Andision/RoundWorld/framework/web/handlers"
	"github.com/Andision/RoundWorld/framework/web/structures"
	"github.com/Andision/RoundWorld/games"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	lounge := structures.NewLounge()
	gameConfig := games.NewConfigGames()
	ctx := context.Background()

	r := mux.NewRouter()

	// 注册 API 接口
	r.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/api/table", structures.JwtValidator(ctx, func(writer http.ResponseWriter, request *http.Request) {
		handlers.CreateTableHandler(request.Context(), lounge, gameConfig, writer, request)
	})).Methods("POST")

	// 注册 WebSocket 接口
	r.HandleFunc("/ws", structures.JwtValidator(ctx, func(w http.ResponseWriter, r *http.Request) {
		log.Println("ws")
		handlers.JoinHandler(r.Context(), lounge, w, r)
	}))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
