// Package api - APIs for front-end
package api

import (
	"encoding/json"
	"net/http"

	"github.com/acornsoft-edgecraft/edgecraft-api/pkg/db"
)

// ===== [ Constants and Variables ] =====

// ===== [ Types ] =====

// Config - Represents the API Server configuration
type Config struct {
	Type        string   `yaml:"type"`
	Port        string   `yaml:"port"`
	Host        string   `yaml:"host"`
	HomePageURL string   `yaml:"homePageUrl"`
	Secret      string   `yaml:"secret"`
	PathPrefix  string   `yaml:"pathPrefix"`
	Langs       []string `yaml:"langs"`
	LangPath    string   `yaml:"langPath"`
	Mode        string   `yaml:"mode"`
	// EdgeDatabase postgresdb.Config `yaml:"edge_database"`
}

// API - Represents the structure of the API
type API struct {
	Config *Config
	Db     db.DB
	// EdgeDb map[string]db.DB
	// Mail   *mail.Client
}

// Middleware - Represents the structure of Middleware
type Middleware struct {
	Name string
}

// Route - Represents the structure of Route
type Route struct {
	Host       string
	RouteRules []RouteRule
}

// RouteRule - Represents the structure of Route Rule
type RouteRule struct {
	Hosts       []string
	Path        string
	Method      string
	Middlewares []Middleware
}

// DataType request data type, JsonType or YamlType
type DataType int

// DataType을 위한 Value Constant
const (
	JSONType = iota //JsonType Json request type
	YAMLType
)

// ===== [ Implementations ] =====

// ===== [ Private Functions ] =====

// getRequestData - Returns the request body data
func getRequestData(req *http.Request, data interface{}) error {
	return json.NewDecoder(req.Body).Decode(data)
}

// ===== [ Public Functions ] =====

// New - Returns the api settings
func New(conf *Config, db db.DB) (*API, error) {
	api := &API{
		Config: conf,
		Db:     db,
	}

	return api, nil
}
