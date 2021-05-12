package checker

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

//go:embed post.html
var indexPage []byte

type Handler interface {
	PutFlag(message *TaskMessage) (*ResultMessage, error)
	GetFlag(message *TaskMessage) (*ResultMessage, error)
	PutNoise(message *TaskMessage) (*ResultMessage, error)
	GetNoise(message *TaskMessage) (*ResultMessage, error)
	Havoc(message *TaskMessage) (*ResultMessage, error)
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
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	res, err := c.checker(request, p)
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

func (c *Checker) checker(request *http.Request, _ httprouter.Params) (*ResultMessage, error) {
	var tm TaskMessage
	if err := json.NewDecoder(request.Body).Decode(&tm); err != nil {
		return nil, err
	}

	switch tm.Method {
	case TaskMessageMethodPutFlag:
		return c.handler.PutFlag(&tm)
	case TaskMessageMethodGetFlag:
		return c.handler.GetFlag(&tm)
	case TaskMessageMethodPutNoise:
		return c.handler.PutNoise(&tm)
	case TaskMessageMethodGetNoise:
		return c.handler.GetNoise(&tm)
	case TaskMessageMethodHavoc:
		return c.handler.Havoc(&tm)
	}

	return nil, fmt.Errorf("method not allowed")
}
