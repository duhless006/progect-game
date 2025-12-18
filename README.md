# 🚀 Угольная компания - Игра с сохранением в PostgreSQL и Redis

## 📋 О проекте

**Угольная компания** - это игра-симулятор управления угледобывающей компанией, написанная на Go с использованием PostgreSQL для сохранения состояния игры и Redis для кэширования и сессий. Проект полностью контейнеризован с помощью Docker. Ваша задача — максимально эффективно управлять ресурсами, чтобы достичь цели в кратчайшие сроки.

## Цель игры
Цель игры
Конечная цель — приобрести всё необходимое оборудование для продуктивной работы предприятия за минимальное время. Для этого вам нужно:

1. Управлять балансом — изначально предприятие имеет пассивный доход 1 уголь в секунду

2. Нанять шахтёров — для увеличения добычи и ускорения накоплений

3. Покупать оборудование — необходимое для завершения игры

## Экономика игры
В игре используется уголь как единая валюта — он одновременно является добываемым материалом и разменной монетой для всех операций.

Шахтёры и их характеристики:
Тип	Стоимость	Энергия	Добыча/раз	Интервал	Прогрессия
Маленький	5 угля	30 добыч	1 уголь	3 сек	Нет
Нормальный	50 угля	45 добыч	3 угля	2 сек	Нет
Сильный	450 угля	60 добыч	10+ угля*	1 сек	+3 с каждой добычей
*Начинает с 10, увеличивается на 3 с каждой добычей

Необходимое оборудование:
*Кирки — 3 000 угля

*Вентиляция — 15 000 угля

*Вагонетки — 50 000 угля

## Стратегия игры
Игроку предстоит принимать стратегические решения:

Когда нанимать шахтёров и каких классов

Как распределять средства между наймом и покупкой оборудования

Как оптимизировать время достижения цели

Каждый шахтёр работает в отдельной горутине и после исчерпания энергии завершает работу. Предприятие может в любой момент завершить работу всех шахтёров.

## Завершение игры
Игра считается завершённой, когда:

Куплено всё оборудование (кирки, вентиляция, вагонетки)

Игрок отправляет команду на завершение

После завершения игра предоставляет статистику:

Общее затраченное время

Количество нанятых шахтёров по классам

Финансовые показатели
## 🏗️ Архитектура проекта
PROGECT-GAME/
├── company/           # Ядро игры: шахтёры, оборудование, бизнес-логика
├── database/          # Работа с данными (PostgreSQL + Redis)
├── http/              # HTTP API сервер, DTO, обработчики запросов
├── models/            # Структуры данных игры
├── .dockerignore      # Исключения для Docker
├── .env               # Переменные окружения
├── docker-compose.yml # Конфигурация Docker (PostgreSQL + Redis + App)
├── dockerfile         # Сборка Go приложения
├── go.mod             # Зависимости Go
├── go.sum             # Контрольные суммы зависимостей
├── init.sql           # Инициализация таблиц PostgreSQL
├── main.go            # Точка входа
├── Makefile           # Автоматизация: тесты, сборка, запуск
└── READMI.md          # Документация проекта

## 🛠️ Технологии

- **Go 1.24.6** - основной язык
- **PostgreSQL 15** - база данных для сохранения игры
- **Redis 7** - кэширование и управление сессиями
- **Docker & Docker Compose** - контейнеризация
- **Gorilla Mux** - HTTP роутинг
- **lib/pq** - драйвер PostgreSQL для Go
- **go-redis** - Redis клиент для Go
- **testify** - фреймворк для тестирования

## 🚀 ПОЭТАПНЫЙ ЗАПУСК ПРОЕКТА

### ✅ ШАГ 0: Предварительные требования

Установите на компьютер:
1. [Docker Desktop](https://www.docker.com/products/docker-desktop/)
2. [Git](https://git-scm.com/)
3. [Go 1.21+](https://go.dev/dl/) (опционально, для разработки)


ШАГ 1: Проверка структуры файлов
Убедись, что в папке есть все необходимые файлы:

bash
ls -la
Должны быть:

main.go

Dockerfile

docker-compose.yml

init.sql

go.mod

Makefile (для тестирования)

ШАГ 2: Первый запуск (сборка и запуск)
bash
# 1. Собери и запусти все контейнеры
docker-compose up --build

# 2. ИЛИ запусти в фоновом режиме
docker-compose up --build -d
Жди 30-60 секунд! PostgreSQL и Redis запускаются дольше Go приложения.

ШАГ 3: Проверка что всё запустилось
bash
# 1. Проверь статус контейнеров
docker-compose ps

# Должно быть:
# NAME                    SERVICE   STATUS        PORTS
# progect-game-db-1       db        running       0.0.0.0:5433->5432/tcp
# progect-game-redis-1    redis     running       0.0.0.0:6379->6379/tcp
# progect-game-app-1      app       running       0.0.0.0:9091->9091/tcp

# 2. Проверь логи
docker-compose logs app | tail -20
# Должно быть: "База данных подключена", "Redis подключен" и "Server starting on :9091"
ШАГ 4: Тестирование API
bash
# 1. Сохранить игру (ручное сохранение)
curl -X POST http://localhost:9091/api/save

# Ответ: {"status":"saved","data":{"money":1000,...}}

# 2. Загрузить игру
curl http://localhost:9091/api/load

🧪 ТЕСТИРОВАНИЕ
Запуск тестов
bash
# Используйте Makefile для удобства
make help                    # Показать все доступные команды

# Все тесты
make test                    # Запустить все тесты с подробным выводом

# Тесты по компонентам
make test-company            # Тесты бизнес-логики компании
make test-db                 # Тесты базы данных
make test-handlers           # Тесты HTTP handlers

# Отчет о покрытии кода
make test-coverage           # Генерирует HTML отчет
# После запуска откройте coverage.html в браузере

# Быстрые юнит-тесты
make test-unit

# Очистка
make clean                   # Удалить сгенерированные файлы тестов
Запуск тестов без Makefile
bash
# Все тесты
go test ./... -v

# Конкретный модуль
go test ./company -v
go test ./database -v
go test ./http -v

# С покрытием кода
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # или xdg-open coverage.html на Linux
Структура тестов
text
tests/
├── company/
│   ├── company_test.go           # Тесты основной логики
│   
├── database/
│   ├── db_test.go                # Тесты PostgreSQL
│   
└── http/
    ├── handlers_test.go          # Тесты API endpoints
   

🔧 КОМАНДЫ ДЛЯ РАБОТЫ С ПРОЕКТОМ
🐳 Docker команды
bash
# Запуск проекта
docker-compose up              # с выводом логов
docker-compose up -d          # в фоновом режиме
docker-compose up --build     # с пересборкой

# Остановка проекта
docker-compose down           # остановить контейнеры
docker-compose down -v        # остановить и удалить тома (данные БД и Redis)

# Просмотр логов
docker-compose logs app       # логи приложения
docker-compose logs db        # логи базы данных
docker-compose logs redis     # логи Redis
docker-compose logs -f app   # логи в реальном времени

# Проверка статуса
docker-compose ps            # статус контейнеров
docker ps                    # все контейнеры Docker

# Пересборка
docker-compose build         # пересобрать образы
docker system prune -f       # очистить Docker (осторожно!)
🗄️ Команды для работы с Redis
bash
# Подключиться к Redis CLI
docker-compose exec redis redis-cli

# Внутри Redis проверить ключи
KEYS *
GET company:default_company
INFO memory

# Или одной командой
docker-compose exec redis redis-cli KEYS "*"
docker-compose exec redis redis-cli GET "company:default_company"
🌐 API эндпоинты
Метод	URL	Описание
POST	http://localhost:9091/api/save	Сохранить текущее состояние игры
GET	http://localhost:9091/api/load	Загрузить последнее сохранение
POST	http://localhost:9091/miners	Нанять нового шахтёра
GET	http://localhost:9091/miners	Получить список шахтёров
POST	http://localhost:9091/equipment	Купить оборудование
GET	http://localhost:9091/company	Получить статистику компании
POST	http://localhost:9091/company/complete	Проверка выполнены ли условия игры
GET	http://localhost:9091/api/cache/status	Статус кэша Redis
DELETE	http://localhost:9091/api/cache/clear	Очистить кэш Redis

## JSON для POST запросов на майнеров и оборудование:
## Miners
-miners little:

{
    "MinerType":"little"
}
-miners normal:

{
    "MinerType":"normal"
}
-miners powerful:

{
    "MinerType":"powerful"
}
## Equpment
-pickaxe - 3000$:

{
    "EquipmentType":"pickaxe"
}
-ventilation - 15000$:

{
    "EquipmentType":"ventilation"
}
-trolleys - 50000$:

{
    "EquipmentType":"trolleys"
}

📊 Структура базы данных
PostgreSQL таблица game_state
id - уникальный идентификатор сохранения

company_id - название компании

money - текущий баланс

total_earned - всего заработано

miners_count - количество шахтёров

pickaxe, ventilation, trolleys - куплено ли оборудование

saved_at - время сохранения

Redis ключи
company:{id} - кэшированное состояние компании

session:{token} - сессия пользователя

leaderboard - таблица лидеров

rate_limit:{ip} - ограничение запросов

🚨 Решение частых проблем
❌ Проблема: "Port 5432 already in use"
Решение: Используй порт 5433 (уже настроено в docker-compose.yml)

❌ Проблема: "Port 6379 already in use" (Redis)
Решение: Остановите локальный Redis или измените порт в docker-compose.yml

❌ Проблема: "relation 'game_state' does not exist"
Решение:

bash
# Создать таблицу вручную
docker-compose exec db psql -U user -d game -c "
CREATE TABLE IF NOT EXISTS game_state (...)"
❌ Проблема: Тесты не запускаются
Решение: Установите зависимости:

bash
go mod download
go mod tidy
go install github.com/stretchr/testify@latest
❌ Проблема: Redis не подключается
Решение: Проверьте что Redis запущен:

bash
docker-compose ps | grep redis
docker-compose logs redis | tail -10
❌ Проблема: Нет покрытия тестами
Решение: Используйте правильную команду:

bash
make test-coverage
# или
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
🔄 Обновление зависимостей
bash
# Обновить все зависимости
go get -u ./...

# Обновить конкретную библиотеку
go get -u github.com/stretchr/testify

# Проверить устаревшие зависимости
go list -u -m all


🆘 Получение помощи
Если что-то не работает:

Проверь логи: docker-compose logs app

Проверь БД: docker-compose exec db psql -U user -d game -c "\dt"

Проверь Redis: docker-compose exec redis redis-cli PING

Запусти тесты: make test для диагностики

Перезапусти: docker-compose down && docker-compose up --build
