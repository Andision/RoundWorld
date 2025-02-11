package main

import (
	"context"
	"github.com/Andision/RoundWorld/framework/web"
	"github.com/Andision/RoundWorld/games"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	lounge := web.NewLounge()
	gameConfig := games.NewConfigGames()
	ctx := context.Background()

	r := mux.NewRouter()

	// 注册 API 接口
	r.HandleFunc("/api/login", web.LoginHandler).Methods("POST")
	r.HandleFunc("/api/create", web.JwtValidator(ctx, func(writer http.ResponseWriter, request *http.Request) {
		web.CreateHandler(request.Context(), lounge, gameConfig, writer, request)
	})).Methods("POST")

	// 注册 WebSocket 接口
	r.HandleFunc("/ws", web.JwtValidator(ctx, func(w http.ResponseWriter, r *http.Request) {
		log.Println("ws")
		web.JoinHandler(r.Context(), lounge, w, r)
	}))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}
