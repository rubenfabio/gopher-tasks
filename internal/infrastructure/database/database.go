package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq" // driver Postgres
)

// Open abre a conexão e faz um Ping para validar.
func Open(driver, dsn string, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (*sql.DB, error) {
    db, err := sql.Open(driver, dsn)
    if err != nil {
        return nil, err
    }

    // Ajustes opcionais de pool
    db.SetMaxOpenConns(maxOpenConns)
    db.SetMaxIdleConns(maxIdleConns)
    db.SetConnMaxLifetime(connMaxLifetime)

    // Testa de fato a conexão
    if err := db.Ping(); err != nil {
        return nil, err
    }
    return db, nil
}
