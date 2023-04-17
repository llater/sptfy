package sptfy

import (
	"encoding/json"
	"github.com/pkg/browser"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	CUE_REDIRECT_ENDPOINT          = "http://cue-server.lelandlater.net:10000/redirect"
	CUE_CLIENT_ID                  = "1fc81450cd414fe38f4a614ebc1e4d67"
	CUE_CLIENT_SCOPE               = `user-read-private user-read-email user-read-private user-top-read user-read-playback-state user-modify-playback-state user-read-currently-playing user-read-recently-played`
	SPOTIFY_AUTHORIZATION_ENDPOINT = "https://accounts.spotify.com/authorize"
	SPOTIFY_ACCESS_TOKEN_ENDPOINT  = "https://accounts.spotify.com/api/token"
	SPOTIFY_API_ENDPOINT           = "https://api.spotify.com/v1"
)

type SpotifyOAuthPkceClient struct {
	HttpClient   http.Client
	PkceVerifier *CodeVerifier
}

type pkceAccessTokenRequest struct {
	GrantType    string `json:"grant_type",default:"authorization_code"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
	ClientId     string `json:"client_id"`
	CodeVerifier string `json:"code_verifier"`
}

type pkceAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

/*
	type spotifyPkceAuthRequest struct {
		ClientId            string `json:"client_id"`
		ResponseType        string `json:"response_type"`
		RedirectUri         string `json:"redirect_uri"`
		State               string `json:"state"`
		Scope               string `json:"scope"`
		CodeChallengeMethod string `default:"S256",json:"code_challenge_method"`
		CodeChallenge       string `json"code_challenge"`
	}
*/

func NewSpotifyOAuthPkceClient() (*SpotifyOAuthPkceClient, error) {
	sClient := SpotifyOAuthPkceClient{}

	// Generate code challenge with helper methods
	verifier, err := verifier()
	challenge := verifier.CodeChallengeS256()

	// Assign verifier to returned client
	sClient.PkceVerifier = verifier

	// Generate state
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	state := strconv.Itoa(seededRand.Int())

	authorizationResponses := make(chan *spotifyAuthorizationResponse)

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			state := r.URL.Query().Get("state")
			if (state != "") && (code != "") {
				reply := &spotifyAuthorizationResponse{
					Code:  code,
					State: state,
				}
				// TODO Check state
				authorizationResponses <- reply
			}
		})
		server := &http.Server{
			Addr:    ":10510",
			Handler: mux,
		}
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	req, err := http.NewRequest("GET", SPOTIFY_AUTHORIZATION_ENDPOINT, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", CUE_CLIENT_ID)
	q.Add("scope", CUE_CLIENT_SCOPE)
	q.Add("redirect_uri", "http://localhost:10510/redirect")
	q.Add("state", state)
	q.Add("code_challenge_method", "S256")
	q.Add("code_challenge", challenge)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println(err)
	}
	browser.OpenURL(req.URL.String())
	authorizationResponse := <-authorizationResponses

	accessTokenResp, err := client.PostForm(SPOTIFY_ACCESS_TOKEN_ENDPOINT, url.Values{
		"client_id":     {CUE_CLIENT_ID},
		"code_verifier": {verifier.Value},
		"grant_type":    {"authorization_code"},
		"code":          {authorizationResponse.Code},
		"redirect_uri":  {"http://localhost:10510/redirect"},
	})
	if err != nil {
		return nil, err
	}
	close(authorizationResponses)

	defer accessTokenResp.Body.Close()
	decoder := json.NewDecoder(accessTokenResp.Body)
	var accessToken accessTokenResponse

	err = decoder.Decode(&accessToken)
	if err != nil {
		return nil, err
	}
	sClient.HttpClient.Transport = oauthTransport{http.Transport{}, accessToken.AccessToken}

	return &sClient, nil
}

func (c *SpotifyOAuthPkceClient) Ping() (spotifyId string, err error) {
	r, err := c.HttpClient.Get(SPOTIFY_API_ENDPOINT + "/me")
	if err != nil {
		return "", err
	}
	log.Println(r.StatusCode)
	if r.StatusCode != http.StatusOK {
		log.Println("Ping to Spotify API failed")
		return "", nil
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	var m spotifyMeResponse
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}
	return m.Id, nil
}
