package util

import (
	"fmt"
	"reflect"
)

// ContextKey defines a context key that will be accessed globally throughout the
// GraphQL app
type ContextKey string

// ContextKeys defines the ContextKey that are available throughout
// the app
type ContextKeys struct {
	User       ContextKey
	Provider   ContextKey
	Dataloader ContextKey
}

// ServerConfig defines the configurations that will be used to
// run the Server
type ServerConfig struct {
	Host     string
	Port     string
	Scheme   string
	GraphQL  GraphQLConfig
	Database DatabaseConfig
	Auth     OpenIDClientConfig
}

// GraphQLConfig defines the configurations for gqlgen to be
// attached on the app
type GraphQLConfig struct {
	Path                string
	PlaygroundPath      string
	EnablePlayground    bool
	EnableIntrospection bool
}

// DatabaseConfig defines the configuration for MongoDB
type DatabaseConfig struct {
	URI         string
	Name        string
	Collections map[string]string
}

// OpenIDClientConfig defines a client that will talk to an OpenID Provider
type OpenIDClientConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	Scopes       []string
}

func (c *ContextKey) String() string {
	return reflect.ValueOf(c).String()
}

func getValidHost(host string) string {
	if host == ":" {
		return "localhost"
	}
	return host
}

// ListenEndpoint builds the endpoint string (host + port)
func (s *ServerConfig) ListenEndpoint() string {
	if s.Port == "80" {
		return s.Host
	}
	if s.Host == ":" {
		return s.Host + s.Port

	}
	return s.Host + ":" + s.Port
}

// Endpoint builds a relative URL
func (s *ServerConfig) Endpoint(path string) string {
	return fmt.Sprintf("/%s", path)
}

// SchemeEndpoint build an absolute URL
func (s *ServerConfig) SchemeEndpoint(path string) string {
	if s.Port == "80" {
		return fmt.Sprintf("%s%s%s", s.Scheme, getValidHost(s.Host), path)
	}
	return fmt.Sprintf("%s%s:%s%s", s.Scheme, getValidHost(s.Host), s.Port, path)
}
