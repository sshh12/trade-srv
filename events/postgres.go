package events

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Database is postgres wrapper
type Database struct {
	db *pg.DB
}

func createTables(db *pg.DB) error {
	models := []interface{}{
		(*Event)(nil),
		(*Symbol)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}

// NewPostgresDatabase creates and connects to a postgres database
func NewPostgresDatabase(user string, password string, addr string, dbName string) (*Database, error) {
	db := pg.Connect(&pg.Options{
		User:     user,
		Password: password,
		Addr:     addr,
		Database: dbName,
	})
	if err := createTables(db); err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

// NewPostgresDatabaseFromURL creates and connects to a postgres database from a URL
func NewPostgresDatabaseFromURL(url string) (*Database, error) {
	opts, err := pg.ParseURL(url)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opts)
	if err := createTables(db); err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

// AddEvent inserts an events into the database
func (database *Database) AddEvent(evt *Event) error {
	_, err := database.db.Model(evt).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}

// AddSymbol inserts a symbol into the database
func (database *Database) AddSymbol(sym *Symbol) error {
	_, err := database.db.Model(sym).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}

// GetSymbols get symbols
func (database *Database) GetSymbols() ([]Symbol, error) {
	var symbols []Symbol
	_, err := database.db.Query(&symbols, `SELECT * FROM symbols`)
	if err != nil {
		return nil, err
	}
	return symbols, nil
}
