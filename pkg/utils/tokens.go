package utils

import (
	"net/http"
)

type AccessTokenTransport struct {
	http.Transport
	AccessToken string
}

func (t AccessTokenTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Add("Authorization", "Bearer "+t.AccessToken)
	return t.Transport.RoundTrip(r)
}
