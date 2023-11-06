package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"sambrlon/config"
	"sambrlon/internal/tokens"
)

type PostgresDB struct {
	db *sqlx.DB
}

func NewPostgresDB(cfg config.DBConfig) (*PostgresDB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{
		db: db,
	}, nil
}

func (pg *PostgresDB) SaveKeys(publicKey, privateKey string) error {
	_, err := pg.db.Exec("INSERT INTO api_keys (public_key, private_key) VALUES ($1, $2)", publicKey, privateKey)
	return err
}

func (pg *PostgresDB) GetPrivateKeyByPublicKey(publicKey string) (string, error) {
	var privateKey string
	err := pg.db.Get(&privateKey, "SELECT private_key FROM api_keys WHERE public_key = $1", publicKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("public key not found")
		}
		return "", err
	}
	return privateKey, nil
}

func (pg *PostgresDB) AddTicker(symbol string, price float64) error {
	_, err := pg.db.Exec("INSERT INTO ticker_info (symbol) VALUES ($1)", symbol)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresDB) AddTickerPrice(symbol string, price float64) error {
	// Выбираем предыдущую цену для символа
	var prevPrice float64
	err := pg.db.Get(&prevPrice, "SELECT price FROM ticker_prices WHERE ticker_id = (SELECT id FROM ticker_info WHERE symbol = $1) ORDER BY date_time DESC LIMIT 1", symbol)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	_, err = pg.db.Exec("INSERT INTO ticker_prices (ticker_id, date_time, price, price_difference) VALUES ((SELECT id FROM ticker_info WHERE symbol = $1), $2, $3, $4)", symbol, time.Now(), price, (price-prevPrice)/prevPrice*100)
	if err != nil {
		return err
	}
	return nil
}

func (pg *PostgresDB) GetTickerData(symbol string, dateFrom, dateTo time.Time) ([]tokens.FetchTickerResponse, error) {
	query := `
        SELECT ti.symbol,
            tp.date_time,
            tp.price
        FROM ticker_prices tp
            JOIN ticker_info ti ON tp.ticker_id = ti.id
        WHERE ti.symbol = $1
            AND tp.date_time >= $2
            AND tp.date_time <= $3
        ORDER BY tp.date_time ASC
    `

	var tickerList []tokens.FetchTickerResponse

	if err := pg.db.Select(&tickerList, query, symbol, dateFrom, dateTo); err != nil {
		return nil, err
	}

	return tickerList, nil
}
