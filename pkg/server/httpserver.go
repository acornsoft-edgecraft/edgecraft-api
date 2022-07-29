package server

import (
	echo "github.com/labstack/echo/v4"
)

// NewHTTPServer Create echo server for Http
func NewHTTPServer() *echo.Echo {
	// Echo instance
	httpserver := echo.New()

	return httpserver
}
