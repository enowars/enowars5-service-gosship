package checker

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

//go:embed post.html
var indexPage []byte

type Handler interface {
	PutFlag(ctx context.Context, message *TaskMessage) (*ResultMessage, error)
	GetFlag(ctx context.Context, message *TaskMessage) (*ResultMessage, error)
	PutNoise(ctx context.Context, message *TaskMessage) (*ResultMessage, error)
	GetNoise(ctx context.Context, message *TaskMessage) (*ResultMessage, error)
	Havoc(ctx context.Context, message *TaskMessage) (*ResultMessage, error)
	GetServiceInfo() *InfoMessage
}

type Checker struct {
	log     *logrus.Logger
	router  *httprouter.Router
	info    *InfoMessage
	handler Handler
}

func NewChecker(log *logrus.Logger, handler Handler) *Checker {
	c := &Checker{
		log:     log,
		router:  httprouter.New(),
		handler: handler,
		info:    handler.GetServiceInfo(),
	}
	c.setupRoutes()
	return c
}

func (c *Checker) checkerWithErrorHandler(writer http.ResponseWriter, request *http.Request, p httprouter.Params) {
	var tm TaskMessage
	if err := json.NewDecoder(request.Body).Decode(&tm); err != nil {
		http.Error(writer, "could not parse body", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx, cancel := context.WithTimeout(request.Context(), time.Duration(tm.Timeout)*time.Millisecond)
	defer cancel()
	res, err := c.checker(ctx, &tm)
	if err != nil {
		c.log.Error(err)
		res = &ResultMessage{
			Result:  ResultError,
			Message: err.Error(),
		}
	}
	if err := json.NewEncoder(writer).Encode(res); err != nil {
		c.log.Error(err)
	}
}

func (c *Checker) setupRoutes() {
	c.router.GET("/", c.index)
	c.router.GET("/service", c.service)
	c.router.POST("/", c.checkerWithErrorHandler)
}

func (c *Checker) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c.log.Printf("%s - %s %s", request.RemoteAddr, request.Method, request.URL.EscapedPath())
	c.router.ServeHTTP(writer, request)
}

func (c *Checker) index(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	writer.Header().Set("Content-Type", "text/html")
	_, err := writer.Write(indexPage)
	if err != nil {
		c.log.Error(err)
	}
}

func (c *Checker) service(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(writer).Encode(c.info)
	if err != nil {
		c.log.Error(err)
	}
}

func (c *Checker) checker(ctx context.Context, tm *TaskMessage) (*ResultMessage, error) {
	switch tm.Method {
	case TaskMessageMethodPutFlag:
		return c.handler.PutFlag(ctx, tm)
	case TaskMessageMethodGetFlag:
		return c.handler.GetFlag(ctx, tm)
	case TaskMessageMethodPutNoise:
		return c.handler.PutNoise(ctx, tm)
	case TaskMessageMethodGetNoise:
		return c.handler.GetNoise(ctx, tm)
	case TaskMessageMethodHavoc:
		return c.handler.Havoc(ctx, tm)
	}

	return nil, fmt.Errorf("method not allowed")
}
