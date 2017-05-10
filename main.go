package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var db *gorm.DB

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type todo struct {
	gorm.Model
	Topic string `json:"topic"`
	Done  bool   `json:"done"`
}

func getTodos(c echo.Context) error {
	var todos []todo
	db = db.Find(&todos)
	if db.Error != nil {
		return db.Error
	}

	return c.JSON(http.StatusOK, todos)
}

func newTodo(c echo.Context) error {
	var t todo

	err := c.Bind(&t)
	if err != nil {
		return err
	}

	db = db.Create(&t)
	if db.Error != nil {
		return db.Error
	}

	return c.JSON(http.StatusOK, t)
}

func delTodo(c echo.Context) error {
	db = db.Delete(todo{}, "id = ?", c.Param("id"))
	if db.Error != nil {
		return db.Error
	}

	return c.String(http.StatusOK, "deleted")
}

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "local.db")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(todo{})
}

func main() {
	// Echo instance
	e := echo.New()
	e.Use(middleware.CORS())

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/todos", getTodos)
	e.POST("/todo", newTodo)
	e.DELETE("/todo/:id", delTodo)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
