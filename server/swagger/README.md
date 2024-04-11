# Swagger 配合 Gin 使用
https://github.com/swaggo/swag/blob/master/README_zh-CN.md

## 1.安装 swag 命令
```go
go install github.com/swaggo/swag/cmd/swag
```

## 2.安装 gin-swagger 命令
```go
go get -u github.com/gin-gonic/gin
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

## 3.编写 swagger 注释
主要要导入生成 docs 包
```go
package main

import (
	_ "yourpath/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gin Swagger 示例 API 文档
// @version 1.0
// @description 这是一个使用 Gin 和 Swagger 生成 API 文档的示例。
// @host localhost:8080
// @BasePath /api/v1
func main() {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/api/v1")
	{
		// 在这里注册你的 API 路径和处理函数
		v1.GET("/ping", ping)
	}

	router.Run(":8080")
}

// 使用 Swagger 注释

// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {object} pongResponse
// @Router /ping [get]

type pongResponse struct {
	Message string `json:"message"`
}

func ping(c *gin.Context) {
	c.JSON(200, pongResponse{Message: "pong"})
}
```

## 4.生成 API 文档
```go
swag init
```

## 5.启动服务运行 API 文档
```go
go run main.go
或者
go generate
```

## 6.swagger 注释规则
https://github.com/swaggo/swag/blob/master/README_zh-CN.md#%E5%A3%B0%E6%98%8E%E5%BC%8F%E6%B3%A8%E9%87%8A%E6%A0%BC%E5%BC%8F
```
下为注释中各个Swagger标签的含义解释：
    @title：API文档的标题。
    @version：API版本。
    @description：API的描述信息。
    @BasePath：API的基础路径，所有路由的前缀。

在每个API函数声明的注释中：
    @Summary：此API的简短总结。
    @Description：API的详细描述。
    @Tags：API的标签，可以用来分类。
    @Accept：API可以处理的MIME类型。
    @Produce：API响应的MIME类型。
    @Param：参数的详细描述，包含参数名、参数位置（query, header, path, cookie, body）、数据类型、是否为必须、描述信息。
    @Success：成功回应的状态码与描述，可以提供返回数据的例子。
    @Failure：失败回应的状态码与描述，可以提供返回数据的例子。
    @Router：路由的路径和HTTP方法。
```
