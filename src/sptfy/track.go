package sptfy

import (
	"net/url"
)

type SptfyTrack struct {
	PlaybackUrl url.URL `json:"playback_url"`
	Name        string  `json:"name"`
	Artists     string  `json:"artists"`
	Album       string  `json:"album"`
	IsPlayable  bool    `json:"is_playable"`
	Id          string  `json:"id"`
	Uri         string  `json:"uri"`
	Href        url.URL `json:"href"`
}

type SpotifyAPITrackSearchResponse struct {
	Tracks struct {
		Href  string `json:"href"`
		Items []struct {
			ExternalUrls struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
			Genres []struct{} `json:"genres"` // need a working example of this
			Href   string     `json:"href"`
			Id     string     `json:"id"`
			Images []struct {
				Height int    `json:"height"`
				Url    string `json:"url"`
				Width  string `json:"width"`
			} `json:"images"`
		} `json:"items"`
		Limit int `json:"limit"`
		// Next int `json:"next"` // I don't know the type
		Offset int `json:"offset"`
		// Previous int `json:"previous"` // Same here
		Total int `json:"total"`
	} `json:"tracks"`
}
