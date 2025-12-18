// database/db.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"progect-game/models"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	connStr := "host=db user=user password=pass dbname=game sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Ошибка открытия БД:", err)
		return err
	}

	// Жду пока БД запустится 30 секунд
	log.Println("Ждём запуска PostgreSQL...")
	for i := 0; i < 30; i++ {
		err = DB.Ping()
		if err == nil {
			log.Println("PostgreSQL подключена!")
			return nil
		}
		log.Printf("Попытка %d: БД ещё не готова...", i+1)
		time.Sleep(1 * time.Second)
	}

	log.Println("Не удалось подключиться к PostgreSQL за 30 секунд")
	return err
}
func SaveState(state models.GameState) error {
	if DB == nil {
		log.Println("БД не подключена, пропускаем сохранение")
		return nil
	}

	query := `
        INSERT INTO game_state 
        (company_id, money, total_earned, miners_count, pickaxe, ventilation, trolleys) 
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := DB.Exec(query,
		state.CompanyID,
		state.Money,
		state.TotalEarned,
		state.MinersCount,
		state.Pickaxe,
		state.Ventilation,
		state.Trolleys,
	)
	return err
}

func SaveMoney(money int) error {
	if DB == nil {
		log.Println("БД не подключена, пропускаем сохранение")
		return nil
	}

	_, err := DB.Exec("INSERT INTO game_state (money) VALUES ($1)", money)
	return err
}

func LoadState() (models.GameState, error) {
	var state models.GameState

	if DB == nil {
		log.Println("БД не подключена, возвращаем начальные значения")
		return models.GameState{
			CompanyID:   "default_company",
			Money:       0,
			TotalEarned: 0,
			MinersCount: 0,
			Pickaxe:     false,
			Ventilation: false,
			Trolleys:    false,
			CreatedAt:   time.Now(),
		}, nil
	}

	query := `
        SELECT company_id, money, total_earned, miners_count, 
               pickaxe, ventilation, trolleys, saved_at
        FROM game_state 
        ORDER BY saved_at DESC 
        LIMIT 1`

	err := DB.QueryRow(query).Scan(
		&state.CompanyID,
		&state.Money,
		&state.TotalEarned,
		&state.MinersCount,
		&state.Pickaxe,
		&state.Ventilation,
		&state.Trolleys,
		&state.CreatedAt,
	)

	if err != nil {
		// Если нет записей, возвращаем состояние по умолчанию
		return models.GameState{
			CompanyID:   "default_company",
			Money:       0,
			TotalEarned: 0,
			MinersCount: 0,
			Pickaxe:     false,
			Ventilation: false,
			Trolleys:    false,
			CreatedAt:   time.Now(),
		}, nil
	}

	return state, nil
}

func CreateTablesIfNotExist() error {
	if DB == nil {
		return fmt.Errorf("БД не подключена")
	}

	query := `
        CREATE TABLE IF NOT EXISTS game_state (
            id SERIAL PRIMARY KEY,
            company_id VARCHAR(100) DEFAULT 'default_company',
            money BIGINT DEFAULT 0,
            total_earned BIGINT DEFAULT 0,
            miners_count INT DEFAULT 0,
            pickaxe BOOLEAN DEFAULT false,
            ventilation BOOLEAN DEFAULT false,
            trolleys BOOLEAN DEFAULT false,
            saved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `

	_, err := DB.Exec(query)
	if err != nil {
		log.Printf("Не удалось создать таблицу: %v", err)
		return err
	}

	log.Println("Таблица game_state создана/проверена")
	return nil
}
