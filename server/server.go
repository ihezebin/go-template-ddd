package server

import (
	"context"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/ihezebin/oneness/httpserver"
	"github.com/ihezebin/oneness/httpserver/middleware"
	"github.com/ihezebin/oneness/runner"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ihezebin/go-template-ddd/server/handler"
	_ "github.com/ihezebin/go-template-ddd/server/swagger/docs"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewServer server
// @title Go Template DDD 示例 API 文档
// @version 1.0
// @description 这是一个使用 Gin 和 Swagger 生成 API 文档的示例。
// @host localhost:8080
// @BasePath /
func NewServer(ctx context.Context, port uint) runner.Task {
	handler := httpserver.NewServerHandlerWithOptions(
		httpserver.WithMiddlewares(middleware.Recovery(), Cors()),
		httpserver.WithLoggingRequest(false),
		httpserver.WithLoggingResponse(false),
		httpserver.WithRouters("",
			handler.NewExampleHandler(),
			// ... other handlers
		),
	)

	pprof.Register(handler)
	handler.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return httpserver.NewServer(httpserver.WithHandler(handler), httpserver.WithPort(port))
}
