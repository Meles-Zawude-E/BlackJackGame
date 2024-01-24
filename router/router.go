package router

import (
	"github.com/labstack/echo/v4"
	"github.com/meles-zawude-e/shuffleGame/handler"
)

func SetUpRouter(e *echo.Echo) {
	v1 := e.Group("/api")

	v1.POST("/deck", handler.CreateDeck)
	v1.GET("/deck/:deckID", handler.OpenDeck)
	v1.POST("/deck/draw/:deckID", handler.DrawCard)
}
