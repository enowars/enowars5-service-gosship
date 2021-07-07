package checker

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

//go:embed post.html
var indexPage []byte

type HandlerInfo struct {
	AttackInfo string
	Flag       string
}

func NewExploitInfo(flag string) *HandlerInfo {
	return &HandlerInfo{Flag: flag}
}

func NewPutFlagInfo(attackInfo string) *HandlerInfo {
	return &HandlerInfo{AttackInfo: attackInfo}
}

type Handler interface {
	PutFlag(ctx context.Context, message *TaskMessage) (*HandlerInfo, error)
	GetFlag(ctx context.Context, message *TaskMessage) error
	PutNoise(ctx context.Context, message *TaskMessage) error
	GetNoise(ctx context.Context, message *TaskMessage) error
	Havoc(ctx context.Context, message *TaskMessage) error
	Exploit(ctx context.Context, message *TaskMessage) (*HandlerInfo, error)
	GetServiceInfo() *InfoMessage
}

type MumbleError interface {
	error
	Mumble() bool
}

func NewMumbleError(msg error) MumbleError {
	return mumbleErrorMsg{msg}
}

type mumbleErrorMsg struct {
	error
}

func (m mumbleErrorMsg) Mumble() bool { return true }

var ErrFlagNotFound = NewMumbleError(errors.New("flag not found"))
var ErrNoiseNotFound = NewMumbleError(errors.New("flag not found"))
var ErrVariantIdOutOfRange = errors.New("variantId out of range")

// Logger is a typical logger interface
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type Checker struct {
	log     Logger
	router  *httprouter.Router
	info    *InfoMessage
	handler Handler
}

func NewChecker(log Logger, handler Handler) *Checker {
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

	startTs := time.Now()
	c.log.Infof("[%s] %s - %s", tm.TaskChainId, tm.Method, tm.TeamName)

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	// reduce timeout by 1000ms so we have enough time to close all connections
	ctx, cancel := context.WithTimeout(request.Context(), time.Duration(tm.Timeout-1000)*time.Millisecond)
	defer cancel()

	var res *ResultMessage
	hi, err := c.checker(ctx, &tm)
	if err != nil {
		c.log.Errorf("[%s] %s - %s", tm.TaskChainId, tm.Method, err.Error())
		if err == context.DeadlineExceeded {
			res = NewResultMessageOffline("timeout")
		} else if _, ok := err.(net.Error); ok {
			res = NewResultMessageOffline("network error")
		} else if _, ok := err.(MumbleError); ok {
			res = NewResultMessageMumble(err.Error())
		} else {
			res = NewResultMessageError(err.Error())
		}
	} else {
		res = NewResultMessageOk()
		if hi != nil {
			res.AttackInfo = hi.AttackInfo
			res.Flag = hi.Flag
		}
	}

	c.log.Infof("[%s] %s - done [%dms]", tm.TaskChainId, tm.Method, time.Since(startTs).Milliseconds())
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
	c.log.Infof("%s - %s %s", request.RemoteAddr, request.Method, request.URL.EscapedPath())
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

func (c *Checker) validateVariantId(tm *TaskMessage) error {
	variants := uint64(0)
	switch tm.Method {
	case TaskMessageMethodPutFlag, TaskMessageMethodGetFlag:
		variants = c.info.FlagVariants
	case TaskMessageMethodPutNoise, TaskMessageMethodGetNoise:
		variants = c.info.NoiseVariants
	case TaskMessageMethodHavoc:
		variants = c.info.HavocVariants
	case TaskMessageMethodExploit:
		variants = c.info.ExploitVariants
	}
	if tm.VariantId >= variants {
		return ErrVariantIdOutOfRange
	}
	return nil
}

func (c *Checker) checker(ctx context.Context, tm *TaskMessage) (*HandlerInfo, error) {
	if err := c.validateVariantId(tm); err != nil {
		return nil, err
	}
	switch tm.Method {
	case TaskMessageMethodPutFlag:
		return c.handler.PutFlag(ctx, tm)
	case TaskMessageMethodGetFlag:
		return nil, c.handler.GetFlag(ctx, tm)
	case TaskMessageMethodPutNoise:
		return nil, c.handler.PutNoise(ctx, tm)
	case TaskMessageMethodGetNoise:
		return nil, c.handler.GetNoise(ctx, tm)
	case TaskMessageMethodHavoc:
		return nil, c.handler.Havoc(ctx, tm)
	case TaskMessageMethodExploit:
		return c.handler.Exploit(ctx, tm)
	}

	return nil, fmt.Errorf("method not allowed")
}
