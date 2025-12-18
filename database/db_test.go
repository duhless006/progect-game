// database/db_logic_test.go
package database

import (
	"progect-game/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestModels тестирует модели данных
func TestModels(t *testing.T) {
	t.Run("GameState structure", func(t *testing.T) {
		now := time.Now()
		state := models.GameState{
			CompanyID:   "test_company",
			Money:       1000,
			TotalEarned: 500,
			MinersCount: 2,
			Pickaxe:     true,
			Ventilation: false,
			Trolleys:    false,
			CreatedAt:   now,
		}

		assert.Equal(t, "test_company", state.CompanyID)
		assert.Equal(t, int64(1000), state.Money)
		assert.Equal(t, int64(500), state.TotalEarned)
		assert.Equal(t, 2, state.MinersCount)
		assert.True(t, state.Pickaxe)
		assert.False(t, state.Ventilation)
		assert.False(t, state.Trolleys)
		assert.Equal(t, now, state.CreatedAt)
	})

	t.Run("Default GameState values", func(t *testing.T) {
		var state models.GameState

		assert.Empty(t, state.CompanyID)
		assert.Equal(t, int64(0), state.Money)
		assert.Equal(t, int64(0), state.TotalEarned)
		assert.Equal(t, 0, state.MinersCount)
		assert.False(t, state.Pickaxe)
		assert.False(t, state.Ventilation)
		assert.False(t, state.Trolleys)
	})
}
