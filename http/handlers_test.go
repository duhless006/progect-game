package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"progect-game/company"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHandleGetCompanyStatistics тестирует API статистики
func TestHandleGetCompanyStatistics(t *testing.T) {
	t.Run("Returns valid JSON response", func(t *testing.T) {
		// Создаём тестовую компанию
		ctx := context.Background()
		company := company.NewCompany(ctx)

		// Создаём handlers
		handlers := NewHTTPHeandlers(company)

		// Создаём тестовый HTTP запрос
		req, err := http.NewRequest("GET", "/company", nil)
		require.NoError(t, err)

		// Создаём ResponseRecorder для захвата ответа
		rr := httptest.NewRecorder()

		// Вызываем обработчик
		handler := http.HandlerFunc(handlers.HandleGetCompanyStatistics)
		handler.ServeHTTP(rr, req)

		// Проверяем статус код
		assert.Equal(t, http.StatusOK, rr.Code,
			"HTTP статус должен быть 200 OK")

		// Проверяем Content-Type
		assert.Contains(t, rr.Header().Get("Content-Type"), "application/json",
			"Content-Type должен быть application/json")

		// Проверяем что тело ответа содержит JSON
		body := rr.Body.String()
		assert.Contains(t, body, "balance",
			"Ответ должен содержать поле 'balance'")
		assert.Contains(t, body, "miners_count",
			"Ответ должен содержать поле 'miners_count'")
	})
}
