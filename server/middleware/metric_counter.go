package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// 创建请求计数器中间件
func MetricRequestCounter(meter metric.Meter) (gin.HandlerFunc, error) {
	counter, err := meter.Int64Counter("requests_total")
	if err != nil {
		return nil, errors.Wrap(err, "create request counter failed")
	}

	return func(c *gin.Context) {
		// 记录请求
		counter.Add(c.Request.Context(), 1, metric.WithAttributes(
			attribute.String("path", c.FullPath()),
			attribute.String("method", c.Request.Method),
			attribute.String("status", fmt.Sprintf("%d", c.Writer.Status())),
		))
		c.Next()
	}, nil
}
