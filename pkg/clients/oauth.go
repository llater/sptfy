package clients

import (
	"encoding/json"
	"errors"
	"github.com/llater/sptfy/pkg/models"
	"github.com/llater/sptfy/pkg/utils"
	"github.com/pkg/browser"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// OAuth2
	CLIENT_ID_ENVVAR_NAME     = "SPOTIFY_CLIENT_ID"
	REFRESH_TOKEN_ENVVAR_NAME = "SPOTIFY_REFRESH_TOKEN"
	REDIRECT_URI_ENVVAR_NAME  = "CUEAPI_OAUTH2_REDIRECT_URI"
)

type SpotifyApiClient struct {
	HttpClient        http.Client
	SpotifyApiBaseUrl string
}

func NewSpotifyApiClient(clientId, clientSecret, redirectUri string) (*SpotifyApiClient, error) {
	authorizationEndpoint := "https://accounts.spotify.com/authorize"
	accessTokenEndpoint := "https://accounts.spotify.com/api/token"
	apiEndpoint := "https://api.spotify.com/v1"

	scope := `user-read-private user-read-email user-read-private user-top-read user-read-playback-state user-modify-playback-state user-read-currently-playing user-read-recently-played`
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	state := strconv.Itoa(seededRand.Int())

	authResponses := make(chan *utils.SpotifyAuthorizationResponse)

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
				authResponses <- reply
			}
		})
		server := &http.Server{
			Addr:    ":10010",
			Handler: mux,
		}
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	time.Sleep(100 * time.Millisecond)

	res, err := http.Get("http://localhost:10010/redirect")
	if err != nil {
		log.Println(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Println("Error, not OK")
	}
	req, err := http.NewRequest("GET", authorizationEndpoint, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", clientId)
	q.Add("scope", scope)
	q.Add("redirect_uri", redirectUri)
	q.Add("show_dialog", "false")
	q.Add("state", state)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println(err)
	}

	browser.OpenURL(req.URL.String())

	// Get authorization code from the channel
	rPair := <-authResponses
	resp, err := client.PostForm(accessTokenEndpoint, url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"code":          {rPair.Code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectUri},
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var accessToken utils.SpotifyAccessTokenResponse

	err = decoder.Decode(&accessToken)
	if err != nil {
		return nil, err
	}

	sClient := SpotifyApiClient{}
	sClient.HttpClient.Transport = utils.AccessTokenTransport{http.Transport{}, accessToken.AccessToken}
	sClient.SpotifyApiBaseUrl = apiEndpoint

	statusRequest, err := http.NewRequest("GET", sClient.SpotifyApiBaseUrl+"/me", nil)
	if err != nil {
		return nil, err
	}

	statusResponse, err := sClient.HttpClient.Do(statusRequest)
	if err != nil {
		return nil, err
	}
	if statusResponse.StatusCode != http.StatusOK {
		log.Fatal("Response wrong from Spotify API")
	}
	defer statusResponse.Body.Close()
	body, err := ioutil.ReadAll(statusResponse.Body)
	if err != nil {
		return nil, err
	}

	var meBody utils.SpotifyMeResponse
	if err := json.Unmarshal(body, &meBody); err != nil {
		return nil, err
	}
	log.Println(meBody)

	return &sClient, nil
}

func (c *SpotifyApiClient) Ping() (spotifyId string, err error) {
	r, err := c.HttpClient.Get(c.SpotifyApiBaseUrl + "/me")
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
	var m utils.SpotifyMeResponse
	if err := json.Unmarshal(b, &m); err != nil {
		return "", err
	}
	log.Println(m)
	return m.Id, nil
}

// Hits the Spotify API and returns a SptfyUser, if they exist.
func (c *SpotifyApiClient) GetUserById(id string) (exists bool, user *models.SptfyUser, err error) {
	// Sanity check input
	if len(id) > 96 || len(id) < 6 {
		return false, nil, errors.New("input user id is invalid length")
	}
	r, err := regexp.Compile(`[[:alnum:]]`)
	if err != nil {
		return false, nil, err
	}
	if !r.MatchString(id) {
		return false, nil, errors.New("input errors do not a match Spotify if regex")
	}

	// Get the user from the Spotify API
	sb := strings.Builder{}
	sb.WriteString(c.SpotifyApiBaseUrl)
	sb.WriteString("/users/")
	sb.WriteString(id)
	req, err := http.NewRequest("GET", sb.String(), nil)
	if err != nil {
		return false, nil, err
	}
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return false, nil, err

	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	var userApiResponse utils.SpotifyUserResponse

	if err := decoder.Decode(&userApiResponse); err != nil {
		return false, nil, err
	}

	return true, &models.SptfyUser{
		DisplayName: userApiResponse.DisplayName,
		Id:          userApiResponse.Id,
	}, nil
}
