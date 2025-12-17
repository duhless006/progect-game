// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"progect-game/company"
	"progect-game/database"
	"progect-game/http"
)

func main() {
	log.Println("Подключаемся к базе данных...")
	if err := database.Init(); err != nil {
		log.Println("Не удалось подключиться к БД. Играем без сохранений.")
		log.Println("Ошибка:", err)
	} else {
		defer database.DB.Close()
		log.Println("База данных подключена")
	}

	if err := database.Init(); err != nil {
		log.Println("Не удалось подключиться к БД")
	} else {
		defer database.DB.Close()

		if err := database.CreateTablesIfNotExist(); err != nil {
			log.Println("Не удалось создать таблицу:", err)
		} else {
			log.Println("Таблица создана")
		}
	}

	state, err := database.LoadState()
	if err == nil {
		log.Printf("Загружено сохранение: деньги=%d, шахтёры=%d",
			state.Money, state.MinersCount)
	} else {
		log.Println("Начинаем новую игру")
	}

	if err := database.InitRedis(); err != nil {
		fmt.Println("Redis:", err)
	}
	fmt.Println("Redis подключен")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	comp := company.NewCompany(ctx)
	handlers := http.NewHTTPHeandlers(comp)
	server := http.NewHTTPServer(handlers)

	fmt.Println("Start Game, Server :9091")
	if err := server.Run(); err != nil {
		fmt.Println("error while running HTTP server:", err)
	}
}
