package clients

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/llater/sptfy/pkg/models"
	"github.com/llater/sptfy/pkg/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type SpotifyClientCredentialsClient struct {
	http.Client
}

func NewSpotifyClientCredentialsClient(clientId, clientSecret string) (*SpotifyClientCredentialsClient, error) {
	credentialsClient := &SpotifyClientCredentialsClient{}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, SPOTIFY_ACCESS_TOKEN_ENDPOINT, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	b := []byte{}
	creds := fmt.Appendf(b, "%s:%s", clientId, clientSecret)
	encodedCredentials := base64.StdEncoding.EncodeToString(creds)

	req.Header.Add("Authorization", "Basic "+encodedCredentials)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var token utils.SpotifyAccessTokenResponse

	err = decoder.Decode(&token)
	if err != nil {
		return nil, err
	}
	credentialsClient.Transport = utils.AccessTokenTransport{http.Transport{}, token.AccessToken}

	return credentialsClient, nil

}

func (c *SpotifyClientCredentialsClient) Me() (*models.SptfyUser, error) {
	r, err := c.Get(SPOTIFY_API_ENDPOINT + "/me")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var m utils.SpotifyMeResponse
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		log.Printf("Observed status code: %d\nAPI status code: %d\nMessage: %s", r.StatusCode, m.Error.Status, m.Error.Message)
		return nil, errors.New("/me endpoint did not return 200")
	}
	return &models.SptfyUser{
		DisplayName: m.Name,
		Email:       m.Email,
		Id:          m.Id,
		Href:        m.URLs.SpotifyLink}, nil
}

func (c *SpotifyClientCredentialsClient) Search(q string) (results *utils.SpotifySearchResponse, err error) {
	r, err := c.Get(SPOTIFY_API_ENDPOINT + "/search?type=track&q=" + q)
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
