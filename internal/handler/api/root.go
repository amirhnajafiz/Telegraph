package api

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

type Root struct {
	Logger *zap.Logger
}

func (root Root) Handle(c echo.Context) error {
	c.Set("Access-Control-Allow-Origin", "*")
	return c.JSON(http.StatusOK, "Welcome home")
}

func (root Root) Register(g *echo.Group) {
	g.GET("/", root.Handle)
}
