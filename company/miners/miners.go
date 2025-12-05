package miners

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Coal int64
type MinerClass string

type MinerInfo struct {
	ID         uuid.UUID     `json:"id"`     // Уникальный номер
	MinerClass MinerClass    `json:"class"`  // Класс (little/normal/powerful)
	Energy     int64         `json:"energy"` // Осталось энергии
	Power      int64         `json:"power"`  // Сколько добывает за удар
	Interval   time.Duration `json:"interval"`
}

type Miner interface {
	Run(ctx context.Context) <-chan Coal
	Info() MinerInfo
}
