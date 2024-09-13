package app

import repository "todo/internal/infra"

type TodoComandService struct {
	repository repository.TodoRepositorier
}

func NewTodoCommandService(repository repository.TodoRepositorier) *TodoComandService {
	return &TodoComandService{repository: repository}
}

