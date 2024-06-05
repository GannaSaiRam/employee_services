package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/GannaSaiRam/employee_services/employee"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      employee.Storage
}

func StartServer(listenAddr string, store employee.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/employee", makeHTTPHandlerFunc(s.handleEmployee))
	router.HandleFunc("/employee/{id}", makeHTTPHandlerFunc(s.handleEmployee))
	router.HandleFunc("/employees", makeHTTPHandlerFunc(s.handleAllEmployees))

	log.Println("JSON Api server running on port:", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleEmployee(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetEmployee(w, r)
	} else if r.Method == "POST" {
		return s.handleCreateEmployee(w, r)
	} else if r.Method == "DELETE" {
		return s.handleDeleteEmployee(w, r)
	} else if r.Method == "PUT" {
		return s.handleUpdateEmployee(w, r)
	}
	return fmt.Errorf("method doesn't exist: %s", r.Method)
}

func (s *APIServer) handleCreateEmployee(w http.ResponseWriter, r *http.Request) error {
	createEmployeeRequest := new(employee.CreateEmployee)
	if err := json.NewDecoder(r.Body).Decode(createEmployeeRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}
	employee := employee.NewEmployeeCreation(createEmployeeRequest.Name, createEmployeeRequest.Position, createEmployeeRequest.Salary)
	id_, err := s.store.CreateEmployee(employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	employee.ID = id_
	return WriteJSON(w, http.StatusOK, employee)
}

func (s *APIServer) handleGetEmployee(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	i, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}
	employee, err := s.store.GetEmployeeByID(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	if employee.ID == 0 {
		http.Error(w, "Data doesn't exist with given ID", http.StatusBadRequest)
		return err
	}
	return WriteJSON(w, http.StatusOK, employee)
}

func (s *APIServer) handleUpdateEmployee(w http.ResponseWriter, r *http.Request) error {
	updateEmployeeRequest := new(employee.UpdateEmployee)
	if err := json.NewDecoder(r.Body).Decode(updateEmployeeRequest); err != nil {
		return err
	}
	employee := employee.EmployeeUpdation(updateEmployeeRequest.Position, updateEmployeeRequest.Salary)
	id := mux.Vars(r)["id"]
	i, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}
	employee.ID = i
	_, err = s.store.UpdateEmployee(employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return WriteJSON(w, http.StatusOK, "")
}

func (s *APIServer) handleDeleteEmployee(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	i, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}
	err = s.store.DeleteEmployee(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return WriteJSON(w, http.StatusOK, "")
}

func (s *APIServer) handleAllEmployees(w http.ResponseWriter, r *http.Request) error {
	limit := r.URL.Query().Get("limit")
	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}
	offset := r.URL.Query().Get("offset")
	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return err
	}
	employees, err := s.store.GetEmployees(offset_int, limit_int)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return WriteJSON(w, http.StatusOK, employees)
}
