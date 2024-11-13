package test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/savioruz/mikti-task/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ClearAll() {
	ClearTodos()
	ClearUsers()
}

func ClearUsers() {
	err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE").Error
	if err != nil {
		log.Fatalf("failed to clear users: %+v", err)
	}
}

func ClearTodos() {
	err := db.Exec("TRUNCATE TABLE todos RESTART IDENTITY CASCADE").Error
	if err != nil {
		log.Fatalf("failed to clear todos: %+v", err)
	}
}

func CreateTodos(t *testing.T, count int) {
	for i := 0; i < count; i++ {
		todos := &entity.Todo{
			ID:        uuid.NewString(),
			Title:     fmt.Sprintf("title-%d", i),
			Completed: false,
		}
		err := db.Create(todos).Error
		assert.Nil(t, err)
	}
}

func GetFirstUser(t *testing.T) *entity.User {
	user := new(entity.User)
	err := db.First(user).Error
	assert.Nil(t, err)
	return user
}

func GetFirstTodoID(t *testing.T) string {
	todo := new(entity.Todo)
	err := db.First(todo).Error
	assert.Nil(t, err)
	return todo.ID
}
