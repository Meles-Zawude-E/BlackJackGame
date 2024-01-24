package main

import (
	"github.com/labstack/echo/v4"
	"github.com/meles-zawude-e/shuffleGame/database"
	"github.com/meles-zawude-e/shuffleGame/router"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	database.InitDB()
	database.GameAutomigrateDatabase()
	router.SetUpRouter(e)
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Start(":5050")
}
