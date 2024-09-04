package main

import (
	"todo/src/controller"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "todo/docs"

	"github.com/labstack/echo/v4"
)

// docs is generated by Swag CLI, you have to import it.

// echo-swagger middleware

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	c := controller.NewTodoController()

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	apiV1 := e.Group("/api/v1")
	apiV1.GET("/", c.GetHello())


	e.Logger.Fatal(e.Start(":8080"))
}