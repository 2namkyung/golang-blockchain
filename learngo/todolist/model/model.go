package model

import (
	"time"
)

type Todo struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type DBHandler interface {

	// User
	GetTodos(sessionId string) []*Todo
	AddTodo(name, sessionId string) *Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, complete bool) bool

	Close()
}

// init : 패키지가 처음 초기화 될 때 호출된다
func NewDBHandler(filepath string) DBHandler {
	// handler = newMemoryHandler()
	return newSqliteHandler(filepath)
}
