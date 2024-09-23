package view

import (
	"log/slog"
	"net/http"
	"strconv"
	app "todo/internal/app"
	"todo/internal/infra"

	"github.com/labstack/echo/v4"
)

// TodoResponseViewはTodoのレスポンスビューです。
type TodoResponseView struct {
	ID         uint64 `json:"id"`
	Title      string `json:"title"`
	Completed  bool   `json:"completed" default:"false"`
	LastUpdate string `json:"last_update"`
	CreatedAt  string `json:"created_at"`
}

// TodoRequestViewはTodoのリクエストビューです。
type TodoRequestView struct { 
	Title string `json:"title"`
	Completed bool `json:"completed"`
}

// TodoControllerはTodoのコントローラーです。
type TodoController struct {
	todoComandService app.TodoComandService
	todoRepository infra.TodoRepositorier
}

// NewTodoControllerはTodoのコントローラーを生成します。
func NewTodoController(todoComandService app.TodoComandService, todoRepository infra.TodoRepositorier) *TodoController {
	return &TodoController{
		todoComandService: todoComandService,
		todoRepository: todoRepository,
	}
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
// @Success 200 {array} TodoResponseView
// @Router /todos [get]
func(ctrl *TodoController) FindAllTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		todos, err := ctrl.todoComandService.FindAllCommand()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		todoViews := make([]TodoResponseView, 0, len(todos))
		for _, todo := range todos {
			todoView := TodoResponseView{
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
// @Success 200 {object} TodoResponseView
// @Router /todos/{id} [get]
func(ctrl *TodoController) FindTodoByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		reqTodo := app.NewTodoIDData(id)
		todo, err := ctrl.todoComandService.FindByIdCommand(reqTodo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		todoView := TodoResponseView{
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
// @Param todo body TodoRequestView true "Todo"
// @Success 200 {object} string
// @Router /todos [post]
func(ctrl *TodoController) CreateTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		todoRequestView := new(TodoRequestView)
		if err := c.Bind(todoRequestView); err != nil {
			slog.Info("bind error", "err",err)
			return c.JSON(http.StatusBadRequest, err)
		}

		todo := app.NewToDoData(0, todoRequestView.Title, todoRequestView.Completed)
		err := ctrl.todoComandService.CreateTodoCommand(todo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusCreated, "success")
	}
}



// UpdateTodo godoc
// @Summary Update todo
// @Description update todo
// @ID update-todo
// @Accept  json
// @Produce  json
// @Param id path int true "Todo ID"
// @Param todo body TodoRequestView true "Todo"
// @Success 200 {object} string
// @Router /todos/{id} [put]
func(ctrl *TodoController) UpdateTodo() echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		todoRequestView := new(TodoRequestView)
		if err := c.Bind(todoRequestView); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		todo := app.NewToDoData(id, todoRequestView.Title, todoRequestView.Completed)
		err = ctrl.todoComandService.UpdateTodoCommand(todo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, "success")
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

		todo := app.NewTodoIDData(id)
		err = ctrl.todoComandService.DeleteTodoCommand(todo)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, nil)
	}
}