package middleware

import (
	"fmt"
	"net/http"
	"valet/component/constant"

	"github.com/gin-gonic/gin"
	"github.com/ihezebin/jwt"
	"github.com/ihezebin/olympus/httpserver"
	"github.com/ihezebin/olympus/logger"
)

var authPathMap = make(map[string]bool)

func pathKey(method, path string) string {
	return fmt.Sprintf("%s:%s", method, path)
}

func RegisterAuthPath(method, path string) {
	authPathMap[pathKey(method, path)] = true
}

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		pathKey := pathKey(method, path)
		needAuth := authPathMap[pathKey]

		body := &httpserver.Body[any]{}
		tokenStr := c.GetHeader(constant.HeaderKeyToken)
		if tokenStr == "" {
			tokenStr = c.Query(constant.QueryKeyToken)
		}
		if tokenStr == "" {
			if needAuth {
				body.WithErr(httpserver.ErrWithUnAuthorized())
				c.AbortWithStatusJSON(http.StatusUnauthorized, body)
				return
			}

			c.Next()
			return
		}

		ctx := c.Request.Context()
		token, err := jwt.Parse(tokenStr, constant.TokenSecret)
		if err != nil {
			logger.WithError(err).Errorf(ctx, "parse token error, token: %s", tokenStr)
			body.WithErr(httpserver.ErrorWithAuthorizationFailed(err.Error()))
			c.AbortWithStatusJSON(http.StatusUnauthorized, body)
			return
		}

		faked, err := token.Faked()
		if err != nil {
			logger.WithError(err).Errorf(ctx, "faked token error, token: %s", tokenStr)
			body.WithErr(httpserver.ErrorWithInternalServer())
			c.AbortWithStatusJSON(http.StatusInternalServerError, body)
			return
		}

		if faked {
			body.WithErr(httpserver.ErrorWithAuthorizationFailed("伪造的令牌"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, body)
			return
		}

		if token.Expired() {
			body.WithErr(httpserver.ErrorWithAuthorizationFailed("令牌已过期"))
			c.AbortWithStatusJSON(http.StatusUnauthorized, body)
			return
		}

		// 传递账号ID
		c.Request.Header.Set(constant.HeaderKeyUid, token.Payload().Owner)
		c.Next()
	}
}
