package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	valication "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ihezebin/olympus/httpserver"
	"github.com/ihezebin/olympus/logger"
	"github.com/ihezebin/openapi"

	"github.com/ihezebin/go-template-ddd/application/dto"
	"github.com/ihezebin/go-template-ddd/application/service"
)

type ExampleRouter struct {
	logger  logger.Logger
	service *service.ExampleApplicationService
}

func NewExampleRouter() *ExampleRouter {
	return &ExampleRouter{}
}

func (r *ExampleRouter) RegisterRoutes(router httpserver.Router) {
	r.logger = logger.WithField("handler", "example")
	r.service = service.NewExampleApplicationService(r.logger)

	// registry http handler
	example := router.Group("example")
	example.POST("/login", httpserver.NewHandler(r.Login),
		httpserver.WithOpenAPISummary("示例登录"),
		httpserver.WithOpenAPIDescription("通过账号和密码登录，生成 Token"),
		httpserver.WithOpenAPIResponseHeader("Token", openapi.HeaderParam{
			Description: "登录成功后返回的 Token",
			Required:    true,
		}),
	)
	example.POST("/register", httpserver.NewHandler(r.Register),
		httpserver.WithOpenAPISummary("示例注册"),
		httpserver.WithOpenAPIDescription("录入账号、密码和邮箱地址"),
	)

}
func (r *ExampleRouter) Login(c *gin.Context, req dto.ExampleLoginReq) (*dto.ExampleLoginResp, error) {
	ctx := c.Request.Context()
	if err := valication.ValidateStruct(&req,
		valication.Field(&req.Username, valication.Required),
		valication.Field(&req.Password, valication.Required),
	); err != nil {
		r.logger.WithError(err).Errorf(ctx, "validate struct error, req: %v", req)
		return nil, httpserver.ErrorWithBadRequest()
	}

	return r.service.Login(ctx, req)

}

func (r *ExampleRouter) Register(c *gin.Context, req dto.ExampleRegisterReq) (*dto.ExampleRegisterResp, error) {
	ctx := c.Request.Context()

	if err := valication.ValidateStruct(&req,
		valication.Field(&req.Username, valication.Required),
		valication.Field(&req.Password, valication.Required),
		valication.Field(&req.Email, valication.Required),
	); err != nil {
		r.logger.WithError(err).Errorf(ctx, "validate struct error, req: %v", req)
		return nil, httpserver.ErrorWithBadRequest()
	}

	resp, err := r.service.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	c.PureJSON(http.StatusOK, &httpserver.Body[*dto.ExampleRegisterResp]{
		Code: httpserver.CodeOK,
		Data: resp,
	})

	return nil, nil

}
