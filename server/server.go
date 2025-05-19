package server

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/ihezebin/olympus/httpserver"
	"github.com/ihezebin/olympus/httpserver/middleware"
	"github.com/ihezebin/olympus/runner"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"

	"github.com/ihezebin/go-template-ddd/config"
	lm "github.com/ihezebin/go-template-ddd/server/middleware" // local middleware
	"github.com/ihezebin/go-template-ddd/server/router"
)

/*
curl --location 'http://127.0.0.1:8080/health' -i
HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Headers: Content-Type, Content-Length, Token
Access-Control-Allow-Methods: POST, GET, PUT, DELETE, OPTIONS
Access-Control-Expose-Headers: Access-Control-Allow-Headers, Token
Content-Type: text/plain; charset=utf-8
Traceparent: 00-5e64f77760384d153783c96049550881-b7ad6b7169203331-01
Date: Wed, 14 May 2025 11:30:41 GMT
Content-Length: 2

ok

其中这个 Traceparent 的格式：
00-5e64f77760384d153783c96049550881-b7ad6b7169203331-01
00-traceId-spanId-flags

traceId: 5e64f77760384d153783c96049550881
spanId: b7ad6b7169203331
flags: 01

相关文档：https://opentelemetry.io/docs/specs/otel/context/api-propagators/#w3c-trace-context-requirements
*/

func NewServer(ctx context.Context, conf *config.Config) (runner.Task, error) {
	meter := otel.Meter(fmt.Sprintf("%s/metrics", conf.ServiceName))
	metricRequestCounter, err := lm.MetricRequestCounter(meter)
	if err != nil {
		return nil, errors.Wrap(err, "create metric request counter failed")
	}

	// logExporter, err := otlploghttp.New(ctx,
	// 	otlploghttp.WithInsecure(),
	// 	otlploghttp.WithEndpoint("localhost:43188"),
	// )
	// if err != nil {
	// 	return nil, errors.Wrap(err, "create otlptracehttp exporter failed")
	// }

	// traceExporter, err := otlptracehttp.New(ctx,
	// 	otlptracehttp.WithInsecure(),
	// 	otlptracehttp.WithEndpoint("localhost:54318"),
	// )
	// if err != nil {
	// 	log.Fatalf("无法创建 OTLP trace HTTP exporter: %v", err)
	// }

	server, err := httpserver.NewServer(ctx,
		httpserver.WithPort(conf.Port),
		httpserver.WithServiceName(conf.ServiceName),
		httpserver.WithMetrics(),
		httpserver.WithPprof(),
		httpserver.WithMiddlewares(
			middleware.Recovery(),
			lm.Cors(),
			metricRequestCounter,
			middleware.LoggingRequestWithoutHeader(),
			middleware.LoggingResponseWithoutHeader(),
		),
		// httpserver.WithLogProcessor(logExporter),
		// httpserver.WithTraceExporter(traceExporter),
		httpserver.WithOpenAPInfo(openapi3.Info{
			Version:     "1.0",
			Description: "这是一个使用 Gin 和 OpenAPI 生成 API 文档的示例。",
			Contact: &openapi3.Contact{
				Name:  "ihezebin",
				Email: "ihezebin@gmail.com",
			},
		}),
		httpserver.WithOpenAPIServer(
			openapi3.Server{
				URL:         fmt.Sprintf("http://localhost:%d", conf.Port),
				Description: "本地开发环境",
			},
			openapi3.Server{
				URL:         "http://api.hezebin.com/go-template-ddd",
				Description: "线上环境",
			},
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "new server err")
	}

	server.RegisterRoutes(
		router.NewExampleRouter(),
	)

	err = server.RegisterOpenAPIUI("/openapi", httpserver.StoplightUI)
	if err != nil {
		return nil, errors.Wrap(err, "register openapi ui err")
	}
	_ = server.RegisterOpenAPIUI("/redoc", httpserver.RedocUI)
	_ = server.RegisterOpenAPIUI("/rapidoc", httpserver.RapidocUI)
	_ = server.RegisterOpenAPIUI("/swagger", httpserver.SwaggerUI)

	return server, nil
}
