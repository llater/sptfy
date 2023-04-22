package utils

// API request bodies
type SpotifyAuthorizationRequest struct {
	ClientId     string `json:"client_id"`
	ResponseType string `json:"response_type"`
	RedirectUri  string `json:"redirect_uri"`
	State        string `json:"state"`
	Scope        string `json:"scope"`
}

// API response bodies
type SpotifyAuthorizationResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
	State string `json:"state"`
}

type SpotifyMeResponse struct {
	Error struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
	Name  string `json:"display_name"`
	Email string `json:"email"`
	Id    string `json:"id"`
	URLs  struct {
		SpotifyLink string `json:"spoztify"`
	} `json:"external_urls"`
}

type SpotifyAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
}

type SpotifySearchResponse struct {
	Error struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
	Tracks struct {
		Items []struct {
			Name    string `json:"name"`
			Id      string `json:"id"`
			Artists []struct {
				Name string `json:"name"`
				Id   string `json:"id"`
			} `json:"artists"`
			Album string `json:"album"`
		} `json:"items"`
	} `json:"tracks"`
	Artists struct {
		Items struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"items"`
	} `json:"artists"`
	Albums struct {
		Items struct {
			Name    string `json:"name"`
			Id      string `json:"id"`
			Artists []struct {
				Name string `json:"name"`
				Id   string `json:"id"`
			}
		} `json:"items"`
	} `json:"albums"`
}

type SpotifyUserResponse struct {
	Error       string `json:"error,omitempty"`
	DisplayName string `json:"display_name"`
	Id          string `json:"id"`
}
