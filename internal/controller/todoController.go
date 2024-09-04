package controller

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"todo/internal/model"
	"todo/internal/repository"

	"github.com/labstack/echo/v4"
)

type TodoView struct {
	ID         uint64 `json:"id"`
	Title      string `json:"title"`
	Completed  bool   `json:"completed"`
	LastUpdate string `json:"last_update"`
	CreatedAt  string `json:"created_at"`
}

type TodoController struct {
	todoRepository repository.TodoRepositorier
}

func NewTodoController(todoRepository repository.TodoRepositorier) *TodoController {
	c := &TodoController{ todoRepository: todoRepository }
	return c
}

// GetHello godoc
// @Summary Get hello
// @Description get hello
// @ID get-hello
// @Produce  plain
// @Success 200 {string} string "Hello, World!"
// @Router / [get]
func(ctrl *TodoController) GetHello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
	
}

// FindAllTodo godoc
// @Summary Find all todos
// @Description get all todos
// @ID find-all-todos
// @Produce  json
// @Success 200 {array} TodoView
// @Router /todos [get]
func(ctrl *TodoController) FindAllTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		todos, err := ctrl.todoRepository.FindAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		todoViews := make([]TodoView, 0, len(todos))
		for _, todo := range todos {
			todoView := TodoView{
				ID:         uint64(todo.ID),
				Title:      todo.Title.AsGoString(),
				Completed:  todo.Completed.AsGoBool(),
				LastUpdate: todo.LastUpdate.AsGoString(),
				CreatedAt:  todo.CreatedAt.AsGoString(),
			}
			todoViews = append(todoViews, todoView)
		}
		return c.JSON(http.StatusOK, todoViews)
	}
}

// FindTodoByID godoc
// @Summary Find todo by ID
// @Description get todo by ID
// @ID find-todo-by-id
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 200 {object} TodoView
// @Router /todos/{id} [get]
func(ctrl *TodoController) FindTodoByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		todo, err := ctrl.todoRepository.FindByID(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		todoView := TodoView{
			ID:         uint64(todo.ID),
			Title:      todo.Title.AsGoString(),
			Completed:  todo.Completed.AsGoBool(),
			LastUpdate: todo.LastUpdate.AsGoString(),
			CreatedAt:  todo.CreatedAt.AsGoString(),
		}
		return c.JSON(http.StatusOK, todoView)
	}
}

// CreateTodo godoc
// @Summary Create todo
// @Description create todo
// @ID create-todo
// @Accept  json
// @Produce  json
// @Param todo body TodoView true "Todo"
// @Success 200 {object} TodoView
// @Router /todos [post]
func(ctrl *TodoController) CreateTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		todoView := new(TodoView)
		if err := c.Bind(todoView); err != nil {
			slog.Info("bind error", "err",err)
			return c.JSON(http.StatusBadRequest, err)
		}

		title, err := model.NewTitle(todoView.Title)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		todo := model.NewTodo(
			0,
			title,
			model.NewCompleted(todoView.Completed),
			model.NewLastUpdate(model.NewModelTime(time.Now())),
			model.NewCreatedAt(model.NewModelTime(time.Now())),
		)

		err = ctrl.todoRepository.Create(todo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		resultTodoView := TodoView{
			ID:         uint64(todo.ID),
			Title:      todo.Title.AsGoString(),
			Completed:  todo.Completed.AsGoBool(),
			LastUpdate: todo.LastUpdate.AsGoString(),
			CreatedAt:  todo.CreatedAt.AsGoString(),
		}
		return c.JSON(http.StatusOK, resultTodoView)
	}
}

// UpdateTodo godoc
// @Summary Update todo
// @Description update todo
// @ID update-todo
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Param todo body TodoView true "Todo"
// @Success 200 {object} TodoView
// @Router /todos/{id} [put]
func(ctrl *TodoController) UpdateTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		todoView := new(TodoView)
		if err := c.Bind(todoView); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		title, err := model.NewTitle(todoView.Title)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		todo := model.NewTodo(
			model.NewID(id),
			title,
			model.NewCompleted(todoView.Completed),
			model.NewLastUpdate(model.NewModelTime(time.Now())),
			model.NewCreatedAt(model.NewModelTime(time.Now())),
		)

		err = ctrl.todoRepository.Update(todo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		resultTodoView := TodoView{
			ID:         uint64(todo.ID),
			Title:      todo.Title.AsGoString(),
			Completed:  todo.Completed.AsGoBool(),
			LastUpdate: todo.LastUpdate.AsGoString(),
			CreatedAt:  todo.CreatedAt.AsGoString(),
		}
		return c.JSON(http.StatusOK, resultTodoView)
	}
}

// DeleteTodo godoc
// @Summary Delete todo
// @Description delete todo
// @ID delete-todo
// @Produce  json
// @Param id path int true "Todo ID"
// @Success 200
// @Router /todos/{id} [delete]
func(ctrl *TodoController) DeleteTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		todo := model.NewTodo(
			model.NewID(id),
			nil,
			nil,
			nil,
			nil,
		)

		err = ctrl.todoRepository.Delete(todo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, nil)
	}
}