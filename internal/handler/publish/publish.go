package publish

import (
	"Telegraph/internal/store/message"
	"Telegraph/pkg/validate"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"time"
)

type Publish struct {
	Database *mongo.Database
	Logger   *zap.Logger
}

type Request struct {
	Source string `json:"from"`
	Des    string `json:"to"`
	Msg    string `json:"message"`
}

func (publish Publish) Handle(c echo.Context) error {
	valid := validate.ValidatePublish(c)
	fmt.Println(reflect.TypeOf(valid["validationError"]))
	if valid["validationError"] != nil {
		return c.JSON(http.StatusBadRequest, valid)
	}

	req := new(Request)

	if err := c.Bind(req); err != nil {
		return err
	}

	item := &message.Message{
		From: req.Source,
		To:   req.Des,
		Msg:  req.Msg,
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err := message.Store(publish.Database, ctx, item)
	if err != nil {
		publish.Logger.Error("insert into database failed", zap.Error(err))
		return err
	}

	// TODO 2: Send the message to the destination
	// TODO 3: Notify the destination

	return c.JSON(http.StatusOK, req)
}

func (publish Publish) Register(g *echo.Group) {
	g.POST("/publish", publish.Handle)
}
