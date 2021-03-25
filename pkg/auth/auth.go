package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// ValidTokenType defines the valid prefix a token
	// must include in the Authorization Header of the
	// request. IE: Authorization: Bearer ...
	ValidTokenType = "Bearer"
	// AuthHeaderKey defines the Header Key that the token
	// will be retrieved from
	AuthHeaderKey = "Authorization"
)

var (
	// ErrUserNotInContext defines an error whereby the auth
	// middleware failed or was not able to retrieve user info
	// from the Auth Server
	ErrUserNotInContext = errors.New("user not found in context")
	// ErrUnauthenticated defines an error where a
	ErrUnauthenticated = errors.New("must be authenticated to access this resource")
)

var client *OpenIDClient

// New creates a new OpenID Client, and sets up important connection details.
func New(clientID string, clientSecret string, openIDAutoDiscoveryURL string, scopes []string) (*OpenIDClient, error) {
	c := &OpenIDClient{
		Scopes:       scopes,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	openIDConfig, err := getOpenIDConfig(openIDAutoDiscoveryURL)
	if err != nil {
		return nil, err
	}
	c.OpenIDConfig = openIDConfig
	client = c
	return c, nil
}

// GetClient retrieves initialized OpendIDClient will return error
// if none have been initialized
func GetClient() (*OpenIDClient, error) {
	if client == nil {
		return nil, errors.New("client has not yet been initialized")
	}
	return client, nil
}

// GetUser uses a given a access token to fetch user info from Auth Provider
func (c *OpenIDClient) GetUser(accessToken string) (*User, error) {
	if accessToken == "" {
		return nil, errors.New("invalid access token provided")
	}
	if c.OpenIDConfig.UserInfoEndpoint == "" {
		return nil, errors.New("client's issuer has no userinfo endpoint")
	}
	url := c.OpenIDConfig.UserInfoEndpoint
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set(AuthHeaderKey, fmt.Sprintf("%s %s", ValidTokenType, accessToken))
	q := req.URL.Query()
	scopes := strings.Join(c.Scopes, " ")
	q.Add("scope", scopes)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"user info endpoint responded with a status of: %d, WWW-Authentica: %s",
			res.StatusCode, res.Header.Get("WWW-Authenticate"),
		)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	user := &User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserFromContext does exactly that
func GetUserFromContext(ctx context.Context, key interface{}) (interface{}, error) {
	user := ctx.Value(key)
	return user, nil
}

// Uses the auto discovery url to populate the rest of the OpenID Client's fields
func getOpenIDConfig(openIDAutoDiscoveryURL string) (*OpenIDConfig, error) {
	httpClient := http.DefaultClient
	res, err := httpClient.Get(openIDAutoDiscoveryURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	openIDConfig := &OpenIDConfig{}
	err = json.Unmarshal(body, openIDConfig)
	if err != nil {
		return nil, err
	}

	return openIDConfig, nil
}
