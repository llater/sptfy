package sptfy

type spotifyAuthorizationResponse struct {
	Code  string `json:"code,omitempty"`
	Error string `json:"error,omitempty"`
	State string `json:"state"`
}
