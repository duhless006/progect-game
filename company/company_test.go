package company

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetMoney тестирует получение баланса
func TestGetMoney(t *testing.T) {
	t.Run("New company has zero money", func(t *testing.T) {
		ctx := context.Background()
		comp := NewCompany(ctx)

		money := comp.GetMoney()
		assert.Equal(t, int64(0), money,
			"Новая компания должна начинать с 0 денег")
	})
}

// TestGetMinerCount тестирует подсчёт шахтёров
func TestGetMinerCount(t *testing.T) {
	t.Run("New company has zero miners", func(t *testing.T) {
		ctx := context.Background()
		comp := NewCompany(ctx)

		count := comp.GetMinerCount()
		assert.Equal(t, 0, count,
			"Новая компания должна начинать с 0 шахтёров")
	})
}

// TestGetState тестирует получение состояния
func TestGetState(t *testing.T) {
	t.Run("GetState returns non-empty structure", func(t *testing.T) {
		ctx := context.Background()
		comp := NewCompany(ctx)

		state := comp.GetState()

		// Проверяем обязательные поля
		assert.NotEmpty(t, state.CompanyID,
			"CompanyID не должен быть пустым")
		assert.Equal(t, "default_company", state.CompanyID,
			"CompanyID должен быть 'default_company'")

		// Проверяем что значения не отрицательные
		assert.GreaterOrEqual(t, state.Money, int64(0))
		assert.GreaterOrEqual(t, state.TotalEarned, int64(0))
		assert.GreaterOrEqual(t, state.MinersCount, 0)

		// Булевы поля должны быть определёнными
		// (хоть false, но не паника)
		_ = state.Pickaxe
		_ = state.Ventilation
		_ = state.Trolleys
		assert.NotNil(t, state.CreatedAt)
	})
}
