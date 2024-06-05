package employee

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPGStore() (*PostgresStore, error) {
	// Any of these two
	// var connStr = "user=postgres dbname=goapi password=postgres sslmode=verify-full"
	var connStr = "postgresql://postgres:postgres@localhost:5433/employee_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{Db: db}, nil
}

func (p *PostgresStore) Init() error {
	p.employees = make(map[int]bool)
	return p.createEmployeeTable()
}

func (p *PostgresStore) createEmployeeTable() error {
	query := `create table if not exists employee (
		id serial primary key,
		name varchar(50),
		position varchar(50),
		salary NUMERIC(10, 2),
		created_at timestamp,
		updated_at timestamp
	)`
	_, err := p.Db.Exec(query)
	return err
}
