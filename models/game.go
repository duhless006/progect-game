// models/game.go
package models

import (
	"encoding/json"
	"time"
)

type GameState struct {
	CompanyID   string    `json:"company_id"`
	Money       int64     `json:"money"`
	TotalEarned int64     `json:"total_earned"`
	MinersCount int       `json:"miners_count"`
	Pickaxe     bool      `json:"pickaxe"`
	Ventilation bool      `json:"ventilation"`
	Trolleys    bool      `json:"trolleys"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToJSON преобразует состояние в JSON строку для хранения
func (s *GameState) ToJSON() string {
	data, _ := json.Marshal(s)
	return string(data)
}

// FromJSON загружает состояние из JSON строки
func (s *GameState) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), s)
}
