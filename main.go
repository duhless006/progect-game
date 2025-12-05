package main

import (
	"context"
	"fmt"
	"progect_game/company"
	"progect_game/http"
)

func main() {
	startCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	comp := company.NewCompany(startCtx)
	handlers := http.NewHTTPHeandlers(comp)
	server := http.NewHTTPServer(handlers)

	fmt.Println("Start Game, Server :9091")
	if err := server.Run(); err != nil {
		fmt.Println("error while running HTTP server:", err)
	}

}
