package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client
var ctx = context.Background()

func InitRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	_, err := Redis.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Ошибка подключения к Redis:", err)
		return err
	}

	fmt.Println("Подключение к Redis успешно")
	return Redis.Ping(context.Background()).Err()
}
