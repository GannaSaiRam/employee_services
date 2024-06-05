package employee

import (
	"fmt"
	"time"
)

func NewEmployeeCreation(name, position string, salary float64) *Employee {
	now_time := time.Now().UTC()
	return &Employee{
		Name:      name,
		Position:  position,
		Salary:    salary,
		CreatedAt: now_time,
		UpdatedAt: now_time,
	}
}

func EmployeeUpdation(position string, salary float64) *Employee {
	now_time := time.Now().UTC()
	return &Employee{
		Position:  position,
		Salary:    salary,
		UpdatedAt: now_time,
	}
}

func (p *PostgresStore) CreateEmployee(e *Employee) (int, error) {
	query := `insert into employee (name, position, salary, created_at, updated_at) values($1, $2, $3, $4, $5) RETURNING id`
	now_time := time.Now()
	var id_ int
	err := p.Db.QueryRow(query, e.Name, e.Position, e.Salary, now_time, now_time).Scan(&id_)
	return id_, err
}

func (p *PostgresStore) GetEmployeeByID(id int) (*Employee, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.employees[id] = true
	defer delete(p.employees, id)
	EmployeeRequest := new(Employee)
	query := fmt.Sprintf(`select id, name, position, salary, created_at, updated_at from employee where id=%d`, id)
	rows, err := p.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&EmployeeRequest.ID,
			&EmployeeRequest.Name,
			&EmployeeRequest.Position,
			&EmployeeRequest.Salary,
			&EmployeeRequest.CreatedAt,
			&EmployeeRequest.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return EmployeeRequest, nil
}

func (p *PostgresStore) UpdateEmployee(emp *Employee) (*Employee, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.employees[emp.ID] = true
	defer delete(p.employees, emp.ID)
	now_time := time.Now()
	query := `update employee
	set position= CASE WHEN $1 != '' THEN $1 ELSE position END,
	salary=CASE WHEN $2 != 0 THEN $2 ELSE salary END,
	updated_at=$3
	where id=$4`
	_, err := p.Db.Exec(query, emp.Position, emp.Salary, now_time, emp.ID)
	return emp, err
}

func (p *PostgresStore) DeleteEmployee(id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.employees[id] = true
	defer delete(p.employees, id)
	query := fmt.Sprintf(`delete from employee where id=%d`, id)
	_, err := p.Db.Exec(query)
	return err
}

func (p *PostgresStore) GetEmployees(offset, limit int) ([]Employee, error) {
	employeesRequest := []Employee{}
	query := fmt.Sprintf(`select id, name, position, salary, created_at, updated_at from employee offset %d limit %d`, offset, limit)
	rows, err := p.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		employee := new(Employee)
		err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Position,
			&employee.Salary,
			&employee.CreatedAt,
			&employee.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		employeesRequest = append(employeesRequest, *employee)
	}
	return employeesRequest, nil
}
