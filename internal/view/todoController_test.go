package view_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"todo/internal/app"
	app_mock "todo/internal/app/mock"
	"todo/internal/domain"
	infra_mock "todo/internal/infra/mock"
	"todo/internal/view"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetHello_成功する(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	mockService := app_mock.NewMockTodoComandService(ctrl)
	repository := infra_mock.NewMockTodoRepositorier(ctrl)
	controller := view.NewTodoController(mockService, repository)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, controller.GetHello()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Hello, World!", rec.Body.String())
	}
}

func TestFindAllTodo_成功する_0件(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockService := app_mock.NewMockTodoComandService(ctrl)
	repository := infra_mock.NewMockTodoRepositorier(ctrl)
	controller := view.NewTodoController(mockService, repository)

	mockService.EXPECT().FindAllCommand().Return([]*domain.Todo{}, nil)
	if assert.NoError(t, controller.FindAllTodo()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "[]\n", rec.Body.String())
	}
}

func TestFindTodoByID_成功する_1件(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/todos/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockService := app_mock.NewMockTodoComandService(ctrl)
	repository := infra_mock.NewMockTodoRepositorier(ctrl)
	controller := view.NewTodoController(mockService, repository)

	id := uint64(1)
	title, err := domain.NewTitle("test")
	assert.Nil(t, err)
	resultTodo := domain.NewTodo(domain.NewID(id), *title, domain.NewCompleted(false), domain.NewLastUpdate(domain.NewDomainTime(time.Now())), domain.NewCreatedAt(domain.NewDomainTime(time.Now())))
	reqTodo := app.NewTodoIDData(id)
	mockService.EXPECT().FindByIdCommand(reqTodo).Return(resultTodo, nil)
	if assert.NoError(t, controller.FindTodoByID()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestCreateTodo_新規作成が成功する(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	e := echo.New()
	todoJSON := `{"title":"New Todo","completed":false}`
	req := httptest.NewRequest(http.MethodPost, "/todos", strings.NewReader(todoJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetRequest(req)

	mockService := app_mock.NewMockTodoComandService(ctrl)
	repository := infra_mock.NewMockTodoRepositorier(ctrl)
	controller := view.NewTodoController(mockService, repository)

	todo := app.NewToDoData(0, "New Todo", false)
	mockService.EXPECT().CreateTodoCommand(todo).Return(nil)
	if assert.NoError(t, controller.CreateTodo()(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestUpdateTodo_更新が成功する(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	e := echo.New()
	todoJSON := `{"title":"Updated Todo","completed":true}`
	req := httptest.NewRequest(http.MethodPut, "/todos/1", strings.NewReader(todoJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockService := app_mock.NewMockTodoComandService(ctrl)
	repository := infra_mock.NewMockTodoRepositorier(ctrl)
	controller := view.NewTodoController(mockService, repository)

	todo := app.NewToDoData(1, "Updated Todo", true)
	mockService.EXPECT().UpdateTodoCommand(todo).Return(nil)
	if assert.NoError(t, controller.UpdateTodo()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "success")
	}
}

func TestDeleteTodo_削除が成功する(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockService := app_mock.NewMockTodoComandService(ctrl)
	repository := infra_mock.NewMockTodoRepositorier(ctrl)
	controller := view.NewTodoController(mockService, repository)

	todo := app.NewTodoIDData(1)
	mockService.EXPECT().DeleteTodoCommand(todo).Return(nil)
	if assert.NoError(t, controller.DeleteTodo()(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}