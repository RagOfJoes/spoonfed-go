package auth

// OpenIDConfig defines the fields that will be fetched from auto discovery url
type OpenIDConfig struct {
	AuthEndpoint     string `json:"authorization_endpoint"`
	TokenEndpoint    string `json:"token_endpoint"`
	UserInfoEndpoint string `json:"userinfo_endpoint"`

	// If OpenID discovery is enabled, the end_session_endpoint field can optionally be provided
	// in the discovery endpoint response according to OpenID spec. See:
	// https://openid.net/specs/openid-connect-session-1_0-17.html#OPMetadata
	EndSessionEndpoint string `json:"end_session_endpoint,omitempty"`
	Issuer             string `json:"issuer"`
}

// OpenIDClient defines a client that will talk to an OpenID Provider
type OpenIDClient struct {
	ClientID     string
	ClientSecret string
	Scopes       []string
	OpenIDConfig *OpenIDConfig
}

// UserProfile defines the basic fields for a
// userinfo response with scope: profile
type UserProfile struct {
	Sub       string `json:"sub"`
	FullName  string `json:"name"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

// User defines the fields for a user info response with
// scope: openid profile
type User struct {
	Sub     string      `json:"sub"`
	Email   string      `json:"email"`
	Profile UserProfile `json:"profile,omitempty"`
}
