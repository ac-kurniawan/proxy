package http

import "github.com/labstack/echo/v4"

type HttpServer struct {
	e *echo.Echo
}

func (h *HttpServer) init() {
	h.e = echo.New()
}

func (h *HttpServer) registerController() {
	h.e.POST("/api/register", HandleProxyPost)
	h.e.GET("/*", HandleProxyGet, NewJwtConfigMiddleware("123213"))
	h.e.POST("/*", HandleProxyPost, NewJwtConfigMiddleware("123213"))
}
