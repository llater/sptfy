package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"strings"
	"time"
)

const (
	pkceVerifierLength = 32
)

// Spotify API responses for OAuth + PLCE enhancement
type PkceAccessTokenRequest struct {
	GrantType    string `json:"grant_type",default:"authorization_code"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
	ClientId     string `json:"client_id"`
	CodeVerifier string `json:"code_verifier"`
}

type PkceAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type CodeVerifier struct {
	Value string
}

func Base64URLEncode(str []byte) string {
	encoded := base64.StdEncoding.EncodeToString(str)
	encoded = strings.Replace(encoded, "+", "-", -1)
	encoded = strings.Replace(encoded, "/", "_", -1)
	encoded = strings.Replace(encoded, "=", "", -1)
	return encoded
}

func Verifier() (*CodeVerifier, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, pkceVerifierLength, pkceVerifierLength)
	for i := 0; i < pkceVerifierLength; i++ {
		b[i] = byte(r.Intn(255))
	}
	return CreateCodeVerifierFromBytes(b)
}

func CreateCodeVerifierFromBytes(b []byte) (*CodeVerifier, error) {
	return &CodeVerifier{
		Value: Base64URLEncode(b),
	}, nil
}

func (v *CodeVerifier) CodeChallengeS256() string {
	h := sha256.New()
	h.Write([]byte(v.Value))
	return Base64URLEncode(h.Sum(nil))
}
