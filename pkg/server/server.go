package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/config"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/logger"
	echo "github.com/labstack/echo/v4"
)

// Instance - Represents an instance of the server
type Instance struct {
	Config     *config.Config
	HTTPServer *echo.Echo
	DB         db.DB
}

// ===== [ Implementations ] =====

// Init - Initialize the server
func (i *Instance) Init() {
	i.HTTPServer = NewHTTPServer()
	logger.Infof("Server initialized...")
}

// Start - Starts the server
func (i *Instance) Start() {

	// Startup the HTTP Server in a way that we can gracefully shut it down again
	err := i.HTTPServer.Start(i.Config.API.Host + ":" + i.Config.API.Port)
	if err != http.ErrServerClosed {
		logger.Errorf("HTTP Server stopped unexpected: %s", err.Error())
		i.Shutdown()
	} else if err != nil {
		logger.Errorf("HTTP Server stoped: %s", err.Error())
	} else {
		logger.Infof("HTTP Server stopped normally")
	}

	logger.Infof("HTTP Server started...++++++")
}

// Shutdown - Stops the server
func (i *Instance) Shutdown() {
	// Shutdown all dependencies
	i.DB.CloseConnection()

	// Shutdown HTTP Server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := i.HTTPServer.Shutdown(ctx)
	if err != nil {
		logger.Errorf("Failed to shutdown HTTP Server gracefully: %s", err.Error())
		os.Exit(1)
	}

	logger.Infof("Shutdown HTTP Server...")
	os.Exit(0)
}

// NewInstance - Returns an new instance of server
func NewInstance(conf *config.Config) *Instance {
	return &Instance{
		Config: conf,
	}
}
