# 🚀 Угольная компания - Игра с сохранением в PostgreSQL

## 📋 О проекте

**Угольная компания** - это игра-симулятор управления угледобывающей компанией, написанная на Go с использованием PostgreSQL для сохранения состояния игры. Проект полностью контейнеризован с помощью Docker.

## 🏗️ Архитектура проекта

```
progect-game/
├── 📁 company/          # Бизнес-логика игры (шахтёры, оборудование)
├── 📁 database/         # Работа с PostgreSQL (сохранение/загрузка)
├── 📁 models/          # Структуры данных (GameState)
├── 📁 http/            # HTTP сервер и API эндпоинты
├── main.go             # Точка входа
├── Dockerfile          # Сборка Go приложения
├── docker-compose.yml  # Запуск приложения + PostgreSQL
├── init.sql           # Инициализация таблицы БД
└── .env               # Переменные окружения
```

## 🛠️ Технологии

- **Go 1.24.6** - основной язык
- **PostgreSQL 15** - база данных
- **Docker & Docker Compose** - контейнеризация
- **Gorilla Mux** - HTTP роутинг
- **lib/pq** - драйвер PostgreSQL для Go

## 🚀 ПОЭТАПНЫЙ ЗАПУСК ПРОЕКТА

### ✅ ШАГ 0: Предварительные требования

Установите на компьютер:
1. [Docker Desktop](https://www.docker.com/products/docker-desktop/)
2. [Git](https://git-scm.com/)

### ✅ ШАГ 1: Клонирование и переход в папку проекта

```bash
# 1. Клонируй репозиторий (если есть)
git clone <url-репозитория>
cd progect-game

# ИЛИ создай папку проекта
mkdir progect-game
cd progect-game
```

### ✅ ШАГ 2: Проверка структуры файлов

Убедись, что в папке есть все необходимые файлы:
```bash
ls -la
```

Должны быть:
- `main.go`
- `Dockerfile`
- `docker-compose.yml`
- `init.sql`
- `go.mod`

### ✅ ШАГ 3: Первый запуск (сборка и запуск)

```bash
# 1. Собери и запусти все контейнеры
docker-compose up --build

# 2. ИЛИ запусти в фоновом режиме
docker-compose up --build -d
```

**Жди 30-60 секунд!** PostgreSQL запускается дольше Go приложения.

### ✅ ШАГ 4: Проверка что всё запустилось

```bash
# 1. Проверь статус контейнеров
docker-compose ps

# Должно быть:
# NAME                    SERVICE   STATUS        PORTS
# progect-game-db-1       db        running       0.0.0.0:5433->5432/tcp
# progect-game-app-1      app       running       0.0.0.0:9091->9091/tcp

# 2. Проверь логи
docker-compose logs app | tail -20
# Должно быть: "База данных подключена" и "Server starting on :9091"
```

### ✅ ШАГ 5: Тестирование API

```bash
# 1. Сохранить игру (ручное сохранение)
curl -X POST http://localhost:9091/api/save

# Ответ: {"status":"saved","data":{"money":1000,...}}

# 2. Загрузить игру
curl http://localhost:9091/api/load


### ✅ ШАГ 6: Проверка базы данных

```bash
# 1. Подключиться к PostgreSQL
docker-compose exec db psql -U user -d game

# 2. Внутри psql проверить таблицу
\dt                    # показать таблицы
\d game_state         # структура таблицы
SELECT * FROM game_state;  # данные в таблице
\q                    # выйти
```

**Или одной командой:**
```bash
docker-compose exec db psql -U user -d game -c "SELECT * FROM game_state ORDER BY saved_at DESC;"
```

## 🔧 КОМАНДЫ ДЛЯ РАБОТЫ С ПРОЕКТОМ

### 🐳 Docker команды

```bash
# Запуск проекта
docker-compose up              # с выводом логов
docker-compose up -d          # в фоновом режиме
docker-compose up --build     # с пересборкой

# Остановка проекта
docker-compose down           # остановить контейнеры
docker-compose down -v        # остановить и удалить тома (данные БД)

# Просмотр логов
docker-compose logs app       # логи приложения
docker-compose logs db        # логи базы данных
docker-compose logs -f app   # логи в реальном времени

# Проверка статуса
docker-compose ps            # статус контейнеров
docker ps                    # все контейнеры Docker

# Пересборка
docker-compose build         # пересобрать образы
docker system prune -f       # очистить Docker (осторожно!)
```

### 🗄️ Команды для работы с PostgreSQL

```bash
# Создать таблицу вручную (если init.sql не сработал)
docker-compose exec db psql -U user -d game -c "
DROP TABLE IF EXISTS game_state;
CREATE TABLE game_state (
    id SERIAL PRIMARY KEY,
    company_id VARCHAR(100) DEFAULT 'default_company',
    money BIGINT DEFAULT 1000,
    total_earned BIGINT DEFAULT 0,
    miners_count INT DEFAULT 0,
    pickaxe BOOLEAN DEFAULT false,
    ventilation BOOLEAN DEFAULT false,
    trolleys BOOLEAN DEFAULT false,
    saved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);"

# Посмотреть структуру таблицы
docker-compose exec db psql -U user -d game -c "\d game_state"

# Посмотреть все сохранения
docker-compose exec db psql -U user -d game -c "
SELECT 
    id,
    money,
    total_earned,
    miners_count,
    pickaxe,
    ventilation,
    trolleys,
    saved_at
FROM game_state ORDER BY saved_at DESC;"

# Удалить все сохранения
docker-compose exec db psql -U user -d game -c "TRUNCATE TABLE game_state;"
```

### 🌐 API эндпоинты

| Метод | URL | Описание |
|-------|-----|----------|
| **POST** | `http://localhost:9091/api/save` | Сохранить текущее состояние игры |
| **GET** | `http://localhost:9091/api/load` | Загрузить последнее сохранение |
| **GET** | `http://localhost:9091/api/saves` | Получить список всех сохранений |
| **POST** | `http://localhost:9091/miners` | Нанять нового шахтёра |
| **GET** | `http://localhost:9091/miners` | Получить список шахтёров |
| **POST** | `http://localhost:9091/equipment` | Купить оборудование |
| **GET** | `http://localhost:9091/company` | Получить статистику компании |
| **POST**|  `http://localhost:9091/company/complete` | Проверка выполнены ли условия игры |

## 📊 Структура базы данных

Таблица `game_state` сохраняет:
- `id` - уникальный идентификатор сохранения
- `company_id` - название компании
- `money` - текущий баланс
- `total_earned` - всего заработано
- `miners_count` - количество шахтёров
- `pickaxe`, `ventilation`, `trolleys` - куплено ли оборудование
- `saved_at` - время сохранения

## Проверка таблицы по завершению игры

# 1. Запустить только PostgreSQL (без приложения)
docker-compose up db -d

# 2. Подождать 20 секунд
sleep 20

# 3. Проверить
docker-compose ps
# Должен быть только db контейнер

# 4. Теперь посмотреть таблицу
docker-compose exec db psql -U user -d game -c "SELECT * FROM game_state;"

# 5. Удаляем данные 
docker-compose down -v 


## 🚨 Решение частых проблем

### ❌ Проблема: "Port 5432 already in use"
**Решение:** Используй порт 5433 (уже настроено в docker-compose.yml)

### ❌ Проблема: "relation 'game_state' does not exist"
**Решение:** 
```bash
# Создать таблицу вручную
docker-compose exec db psql -U user -d game -c "
CREATE TABLE IF NOT EXISTS game_state (...)"
```

### ❌ Проблема: БД не подключается
**Решение:** Подожди 30 секунд (PostgreSQL запускается медленно)

### ❌ Проблема: "service 'db' is not running"
**Решение:** 
```bash
docker-compose up db -d  # запустить только БД
sleep 10                 # подождать
docker-compose ps        # проверить статус
```

## 🎮 Как играть

1. **Запусти проект:** `docker-compose up`
2. **Открой в браузере:** http://localhost:9091
3. **Нанять шахтёров:** через API или интерфейс
4. **Купить оборудование:** для увеличения добычи
5. **Сохранить игру:** нажми кнопку или используй `/api/save`
6. **Выход:** Ctrl+C в терминале (игра сохранится автоматически)

## 📈 Дальнейшее развитие

1. **Добавить фронтенд** - веб-интерфейс для игры
2. **Автосохранение** - каждые N минут
3. **Несколько сохранений** - выбор конкретного сохранения
4. **Миграции БД** - для обновления структуры таблиц
5. **Статистика** - графики и отчёты

## 📝 Примечания

- **Данные сохраняются** между запусками (Docker volume)
- **Для разработки** можно монтировать код в контейнер
- **Логи PostgreSQL** можно смотреть через `docker-compose logs db`
- **Порт 5433** используется чтобы не конфликтовать с локальным PostgreSQL

## 🆘 Получение помощи

Если что-то не работает:
1. Проверь логи: `docker-compose logs app`
2. Проверь БД: `docker-compose exec db psql -U user -d game -c "\dt"`
3. Перезапусти: `docker-compose down && docker-compose up --build`

---

**Удачной игры!** 🎮⛏️💰

> *Проект разработан для изучения Go, PostgreSQL и Docker*