package employee

import (
	"database/sql"
	"sync"
	"time"
)

// Interface any database should implement if database going to be different one
type Storage interface {
	CreateEmployee(*Employee) (int, error)
	GetEmployeeByID(int) (*Employee, error)
	UpdateEmployee(*Employee) (*Employee, error)
	DeleteEmployee(int) error
	GetEmployees(int, int) ([]Employee, error)
}

type Employee struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Position  string    `json:"position"`
	Salary    float64   `json:"salary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateEmployee struct {
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}
type UpdateEmployee struct {
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

type PostgresStore struct {
	Db        *sql.DB
	employees map[int]bool
	mu        sync.RWMutex
}
