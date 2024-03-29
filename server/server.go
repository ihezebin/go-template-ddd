package server

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ihezebin/go-template-ddd/server/handler"
	"github.com/ihezebin/go-template-ddd/server/middleware"
	"github.com/ihezebin/sdk/httpserver"
	"github.com/ihezebin/sdk/httpserver/handler/result"
	middle "github.com/ihezebin/sdk/httpserver/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func NewServer(port int) httpserver.Server {
	server := httpserver.NewServer(
		httpserver.WithPort(port),
		httpserver.WithMiddlewares(
			middle.LoggingRequest(),
			middle.LoggingResponse(),
			middle.Recovery(),
			middleware.Cors(),
		),
	)

	router := server.Kernel()
	pprof.Register(router)

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, result.Json{"status": "ok"}) })

	v1 := router.Group("v1")
	// init handlers
	initHandlers(v1)

	return server
}

type Handler interface {
	Init(gin.IRouter)
}

func initHandlers(router gin.IRouter) {
	handlers := []Handler{
		&handler.TestHandler{},
		// ... other handlers
	}
	for _, hdl := range handlers {
		hdl.Init(router)
	}
}
