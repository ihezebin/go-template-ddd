package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	valication "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/ihezebin/go-template-ddd/application/dto"
	"github.com/ihezebin/go-template-ddd/application/service"
	"github.com/ihezebin/oneness/httpserver"
	"github.com/ihezebin/oneness/logger"
)

type ExampleHandler struct {
	logger  *logger.Entry
	service *service.ExampleApplicationService
}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

func (h *ExampleHandler) Init(router gin.IRouter) {
	h.logger = logger.WithField("handler", "example")
	h.service = service.NewExampleApplicationService(h.logger)

	// registry http handler
	if router != nil {
		example := router.Group("example")
		example.POST("/login", httpserver.NewHandlerFunc(h.Login))
		example.POST("/register", httpserver.NewHandlerFuncEnhanced(h.Register))
	}

}

func (h *ExampleHandler) Login(ctx context.Context, req *dto.ExampleLoginReq) (*dto.ExampleLoginResp, error) {
	if err := valication.ValidateStruct(req,
		valication.Field(&req.Username, valication.Required),
		valication.Field(&req.Password, valication.Required),
	); err != nil {
		h.logger.WithError(err).Error(ctx, "validate struct error, req: %v", req)
		return nil, httpserver.ErrorWithBadRequest()
	}

	return h.service.Login(ctx, req)

}

func (h *ExampleHandler) Register(c *gin.Context, req *dto.ExampleRegisterReq) (*dto.ExampleRegisterResp, error) {
	ctx := c.Request.Context()

	if err := valication.ValidateStruct(req,
		valication.Field(&req.Username, valication.Required),
		valication.Field(&req.Password, valication.Required),
		valication.Field(&req.Email, valication.Required),
	); err != nil {
		h.logger.WithError(err).Error(ctx, "validate struct error, req: %v", req)
		return nil, httpserver.ErrorWithBadRequest()
	}

	resp, err := h.service.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	c.PureJSON(http.StatusOK, &httpserver.Body{
		Code: httpserver.CodeOK,
		Data: resp,
	})

	return nil, nil

}
