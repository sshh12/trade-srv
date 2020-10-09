package events

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Database struct {
	db *pg.DB
}

func createTables(db *pg.DB) error {
	models := []interface{}{
		(*Event)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{IfNotExists: true})
		if err != nil {
			return err
		}
	}
	return nil
}

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

func (database *Database) AddEvent(evt *Event) error {
	_, err := database.db.Model(evt).OnConflict("DO NOTHING").Insert()
	if err != nil {
		return err
	}
	return nil
}
