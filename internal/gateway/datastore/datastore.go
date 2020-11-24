package datastore

import (
	"github.com/jmoiron/sqlx"
	"github.com/rknizzle/faas/internal/models"
)

type Datastore struct {
	db *sqlx.DB
}

func NewDatastore() (Datastore, error) {
	// TODO: hardcoded example -- make this configurable
	conn := "user=postgres dbname=db password=password host=localhost sslmode=disable"

	db, err := sqlx.Open("postgres", conn)
	if err != nil {
		return Datastore{}, err
	}
	return Datastore{db}, nil
}

func (ds Datastore) CreateTable() {
	schema := `
		CREATE TABLE functions (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL, // UNIQUE,
			image TEXT NOT NULL,
			updated_at TIMESTAMP DEFAULT NULL,
			created_at TIMESTAMP DEFAULT NULL
		);
	`

	// TODO: I think mustExec panics so maybe not use that command
	ds.db.MustExec(schema)
}

func (ds Datastore) GetByName(name string) (models.Function, error) {
	fn := models.Function{}
	query := ds.db.Rebind(`SELECT * FROM functions WHERE name = ?`)
	row := ds.db.QueryRowx(query, name)

	err := row.StructScan(&fn)
	if err != nil {
		return models.Function{}, nil
	}
	return fn, nil
}

func (ds Datastore) Create(fn models.Function) error {
	query := ds.db.Rebind(`INSERT INTO functions (name, image, updated_at, created_at) VALUES (?, ?, NOW(), NOW()) RETURNING id`)
	_, err := ds.db.NamedExec(query, fn)
	if err != nil {
		return err
	}

	// TODO: what happens if one already exists with that name?

	return nil
}
