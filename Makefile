.PHONY: test test-unit test-coverage test-company test-db test-handlers clean help

# Цвета для вывода
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m

# Запуск всех тестов
test:
	@echo "${YELLOW}Запуск всех тестов...${NC}"
	go test ./... -v

# Запуск юнит-тестов
test-unit:
	@echo "${YELLOW}Запуск юнит-тестов...${NC}"
	go test ./... -v -short

# Покрытие кода
test-coverage:
	@echo "${YELLOW}Генерация отчета о покрытии...${NC}"
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}Отчет сохранен в coverage.html${NC}"
	@echo "Откройте: open coverage.html  # Mac"
	@echo "         xdg-open coverage.html  # Linux"

# Тесты company
test-company:
	@echo "${YELLOW}Запуск тестов company...${NC}"
	go test ./company -v

# Тесты database  
test-db:
	@echo "${YELLOW}Запуск тестов database...${NC}"
	go test ./database -v

# Тесты handlers
test-handlers:
	@echo "${YELLOW}Запуск тестов handlers...${NC}"
	go test ./http -v

# Очистка
clean:
	@echo "${YELLOW}Очистка...${NC}"
	rm -f coverage.out coverage.html

# Показать доступные команды
help:
	@echo "${YELLOW}Доступные команды:${NC}"
	@echo "${GREEN}make test${NC}          - Запустить все тесты"
	@echo "${GREEN}make test-unit${NC}     - Только юнит-тесты"
	@echo "${GREEN}make test-coverage${NC} - Тесты с отчетом покрытия"
	@echo "${GREEN}make test-company${NC}  - Тесты компании"
	@echo "${GREEN}make test-db${NC}       - Тесты базы данных"
	@echo "${GREEN}make test-handlers${NC} - Тесты HTTP handlers"
	@echo "${GREEN}make clean${NC}         - Очистить сгенерированные файлы"
