package sptfy

type SptfyUser struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Id          string `json:"id"`
	Uri         string `json:"uri"`
	Href        string `json:"href"`
}

type SpotifyAPIUserResponse struct {
	Birthdate    string `json:"birthdate"`
	Country      string `json:"country"`
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	ExternalURLs struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`

	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Href   string `json:"href"`
	Id     string `json:"id"`
	Images []struct {
		Height int    `json:"height"`
		Url    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Product string `json:"product"`
	Type    string `json:"type"`
	Uri     string `json:"uri"`
}
