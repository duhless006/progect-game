-- init.sql - ЗАМЕНИ ВЕСЬ ФАЙЛ НА ЭТОТ:
DROP TABLE IF EXISTS game_state;

CREATE TABLE game_state (
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