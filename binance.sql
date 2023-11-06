    -- Создание таблицы для хранения информации о добавленных тикерах
    CREATE TABLE IF NOT EXISTS ticker_info (
        id SERIAL PRIMARY KEY,
        symbol VARCHAR(10) UNIQUE NOT NULL
    );

    -- Создание таблицы для хранения исторических данных о ценах тикеров
    CREATE TABLE IF NOT EXISTS ticker_prices (
        id SERIAL PRIMARY KEY,
        ticker_id INT NOT NULL,
        date_time TIMESTAMP NOT NULL,
        price double precision NOT NULL,
        price_difference double precision DEFAULT 0,
        FOREIGN KEY (ticker_id) REFERENCES ticker_info
    );

    -- Создание таблицы для хранения API-ключей
    CREATE TABLE IF NOT EXISTS api_keys (
        id SERIAL PRIMARY KEY,
        public_key VARCHAR(255) UNIQUE NOT NULL,
        private_key VARCHAR(255) NOT NULL,
        description TEXT
    );
