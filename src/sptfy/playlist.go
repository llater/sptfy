package sptfy

import (
	"net/url"
)

type SptfyPlaylist struct {
	Name   string             `json:"name"`
	Owner  SptfyUser     `json:"owner"`
	Tracks []SptfyTrack `json:"tracks"`
	Id     string             `json:"id"`
	Uri    string             `json:"uri"`
	Href   url.URL            `json:"href"`
}
