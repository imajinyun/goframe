package middleware

import (
	"net/http"
	"strings"

	"github.com/imajinyun/goframe/gin"
)

const (
	allowOrigin      = "Access-Control-Allow-Origin"
	allowMethods     = "Access-Control-Allow-Methods"
	allowHeaders     = "Access-Control-Allow-Headers"
	allowCredentials = "Access-Control-Allow-Credentials"
	exposeHeaders    = "Access-Control-Expose-Headers"
	requestMethod    = "Access-Control-Request-Method"
	requestHeaders   = "Access-Control-Request-Headers"
	maxAgeHeader     = "Access-Control-Max-Age"
	varyHeader       = "Vary"
	originHeader     = "Origin"
)

var defaultCors = Cors{
	allowOrigins: []string{"*"},
	allowMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodHead,
		http.MethodOptions,
	},
	allowHeaders: []string{
		"Content-Type",
		"Origin",
		"X-CSRF-Token",
		"Authorization",
		"AccessToken",
		"Token",
		"Range",
	},
	exposeHeaders: []string{
		"Content-Length",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
	},
	maxAge: "86400",
}

type Cors struct {
	allowOrigins     []string
	allowMethods     []string
	allowHeaders     []string
	exposeHeaders    []string
	allowCredentials bool
	maxAge           string
}

type CorsOption func(c *Cors)

func NewCors(opts ...CorsOption) *Cors {
	cors := defaultCors
	for _, opt := range opts {
		opt(&cors)
	}

	return &cors
}

func (c *Cors) Func() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add(varyHeader, originHeader)
		if ctx.Request.Method == http.MethodOptions {
			ctx.Writer.Header().Add(varyHeader, requestMethod)
			ctx.Writer.Header().Add(varyHeader, requestHeaders)
		}

		origin := ctx.Request.Header.Get(originHeader)
		if b, allowOrigins := isAllowedOrigin(c.allowOrigins, origin); b {
			c.setCorsHeader(ctx, allowOrigins)
		}

		if ctx.Request.Method == http.MethodOptions {
			ctx.JSON(http.StatusNoContent, nil)
			return
		}

		ctx.Next()
	}
}

func (c *Cors) setCorsHeader(ctx *gin.Context, origins []string) {
	ctx.Writer.Header().Set(allowOrigin, strings.Join(origins, ","))
	ctx.Writer.Header().Set(allowMethods, strings.Join(c.allowMethods, ","))
	ctx.Writer.Header().Set(allowHeaders, strings.Join(c.allowHeaders, ","))
	ctx.Writer.Header().Set(exposeHeaders, strings.Join(c.exposeHeaders, ","))

	if c.maxAge != "0" {
		ctx.Writer.Header().Set(maxAgeHeader, c.maxAge)
	}
}

func WithAllOrigins(allowOrigins []string) CorsOption {
	return func(c *Cors) {
		c.allowOrigins = allowOrigins
	}
}

func WithAllHeaders(allowHeaders []string) CorsOption {
	return func(c *Cors) {
		c.allowHeaders = allowHeaders
	}
}

func WithAllowMethods(allowMethods []string) CorsOption {
	return func(c *Cors) {
		c.allowMethods = allowMethods
	}
}

func WithExposeHeaders(exposeHeaders []string) CorsOption {
	return func(c *Cors) {
		c.exposeHeaders = exposeHeaders
	}
}

func WithAllowCredentials(allowCredentials bool) CorsOption {
	return func(c *Cors) {
		c.allowCredentials = allowCredentials
	}
}

func WithMaxAge(maxAge string) CorsOption {
	return func(c *Cors) {
		c.maxAge = maxAge
	}
}

func isAllowedOrigin(allows []string, origin string) (bool, []string) {
	for _, allow := range allows {
		if allow == "*" {
			return true, []string{"*"}
		}

		if allow == origin {
			return true, []string{origin}
		}
	}

	return false, []string{}
}
