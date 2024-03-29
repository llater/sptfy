package clients

import (
	"encoding/json"
	"fmt"
	"github.com/llater/sptfy/pkg/models"
	"github.com/llater/sptfy/pkg/utils"
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
	CUE_REDIRECT_ENDPOINT = "http://cue-server.lelandlater.net:10000/redirect"
	CUE_CLIENT_ID         = "1fc81450cd414fe38f4a614ebc1e4d67"
	CUE_CLIENT_SCOPE      = `user-read-private user-read-email user-read-private user-top-read user-read-playback-state user-modify-playback-state user-read-currently-playing user-read-recently-played`
)

type SpotifyOAuthPkceClient struct {
	HttpClient   http.Client
	PkceVerifier *utils.CodeVerifier
}

func NewSpotifyOAuthPkceClient() (*SpotifyOAuthPkceClient, error) {
	sClient := SpotifyOAuthPkceClient{}

	// Generate code challenge with helper methods
	verifier, err := utils.Verifier()
	challenge := verifier.CodeChallengeS256()

	// Assign verifier to returned client
	sClient.PkceVerifier = verifier

	// Generate state
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	state := strconv.Itoa(seededRand.Int())

	authorizationResponses := make(chan *utils.SpotifyAuthorizationResponse)

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			state := r.URL.Query().Get("state")
			if (state != "") && (code != "") {
				reply := &utils.SpotifyAuthorizationResponse{
					Code:  code,
					State: state,
				}
				// TODO Check state
				authorizationResponses <- reply
				http.Redirect(w, r, "http://localhost:10011/", http.StatusSeeOther)
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
	var accessToken utils.SpotifyAccessTokenResponse

	err = decoder.Decode(&accessToken)
	if err != nil {
		return nil, err
	}
	sClient.HttpClient.Transport = utils.AccessTokenTransport{http.Transport{}, accessToken.AccessToken}

	return &sClient, nil
}

func (c *SpotifyOAuthPkceClient) Search(q string) (results *utils.SpotifySearchResponse, err error) {
	r, err := c.HttpClient.Get(SPOTIFY_API_ENDPOINT + "/search?type=track&q=" + q)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		log.Println("Failed to reach Spotify API /search endpoint")
		return nil, fmt.Errorf("Failed to reach Spotify API /me endpoint with status code %d", r.StatusCode)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var s utils.SpotifySearchResponse
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (c *SpotifyOAuthPkceClient) Me() (user *models.SptfyUser, err error) {
	r, err := c.HttpClient.Get(SPOTIFY_API_ENDPOINT + "/me")
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		log.Println("Failed to reach Spotify API /me endpoint")
		return nil, fmt.Errorf("Failed to reach Spotify API /me endpoint with status code %d", r.StatusCode)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var m utils.SpotifyMeResponse
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return &models.SptfyUser{
		DisplayName: m.Name,
		Email:       m.Email,
		Id:          m.Id,
		Href:        m.URLs.SpotifyLink}, nil
}

func (c *SpotifyOAuthPkceClient) Ping() error {
	r, err := c.HttpClient.Get(SPOTIFY_API_ENDPOINT + "/me")
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		log.Println("Ping to Spotify API failed")
		return fmt.Errorf("Failed to reach Spotify API /me endpoint with status code %d", r.StatusCode)
	}
	return nil
}
