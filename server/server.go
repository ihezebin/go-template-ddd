package server

import (
	"context"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ihezebin/go-template-ddd/server/handler"
	"github.com/ihezebin/go-template-ddd/server/middleware"
	"github.com/ihezebin/oneness/httpserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(ctx context.Context, port uint) error {
	serverHandler := httpserver.NewServerHandlerWithOptions(
		httpserver.WithLoggingRequest(false),
		httpserver.WithLoggingResponse(false),
		httpserver.WithMiddlewares(middleware.Cors()),
		httpserver.WithRouters("",
			handler.NewExampleHandler(),
			// ... other handlers
		),
	)

	pprof.Register(serverHandler)
	serverHandler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	httpserver.ResetServerHandler(serverHandler)

	return httpserver.Run(ctx, port)
}
