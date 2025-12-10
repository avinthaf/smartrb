package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)


func New(context context.Context, connection_str string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connection_str)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(context); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

var defaultDB *sql.DB

func SetDefault(db *sql.DB) { defaultDB = db }
func Default() *sql.DB      { return defaultDB }

func Seed(db *sql.DB) error {
    // Read the SQL file
    content, err := os.ReadFile("seed.sql")
    if err != nil {
        return err
    }

    // Split by semicolon to separate statements
    statements := strings.Split(string(content), ";")
    
    for _, statement := range statements {
        statement = strings.TrimSpace(statement)
        if statement == "" {
            continue
        }
        if _, err := db.Exec(statement); err != nil {
            return fmt.Errorf("failed to execute statement: %v\nStatement: %s", err, statement)
        }
    }
    
    return nil
}
