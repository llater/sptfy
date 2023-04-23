package main

import (
	"flag"
	"fmt"
	"github.com/llater/sptfy/pkg/clients"
	"github.com/llater/sptfy/pkg/models"
	"github.com/llater/sptfy/pkg/utils"
	"github.com/manifoldco/promptui"
	"log"
	"net/url"
	"os"
)

var (
	spotifyClientId         string
	spotifyClientIdPath     string
	spotifyClientSecret     string
	spotifyClientSecretPath string
)

func crash(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	homedir, err := os.UserHomeDir()
	crash(err)

	// Filepath defaults to the home directory ~/.cue
	// These secrets can only be read in from a filepath.
	flag.StringVar(&spotifyClientIdPath, "spotify-client-id-path", fmt.Sprintf("%s%s", homedir, "/.sptfy/spotify-client-id"), "Spotify API client ID")
	flag.StringVar(&spotifyClientSecretPath, "spotify-client-secret", fmt.Sprintf("%s%s", homedir, "/.sptfy/spotify-client-secret"), "Spotify API client secret")
	flag.Parse()

	// Read in the flags from the provided filepath
	spotifyClientId, err := os.ReadFile(spotifyClientIdPath)
	crash(err)

	spotifyClientSecret, err := os.ReadFile(spotifyClientSecretPath)
	crash(err)

	client, err := clients.NewSpotifyClientCredentialsClient(spotifyClientId, spotifyClientSecret)
	crash(err)

	var command string
	if len(os.Args) < 2 {
		log.Fatal("Display")
	} else {
		command = os.Args[1] // Call only the first argument
	}

	switch command {
	case "search":
		var (
			query    string
			response *utils.SpotifySearchResponse
		)
		if len(os.Args) < 3 {
			queryPrompt := promptui.Prompt{
				Label: "Search",
			}
			query, err = queryPrompt.Run()
			crash(err)
		} else {
			// Use the next argument as the search query
			query = os.Args[2]
		}
		query = url.QueryEscape(query)
		response, err = client.Search(query)
		crash(err)

		tracks := response.Tracks.Items
		outputTracks := []*models.SptfyTrack{}

		for i := 0; i < len(tracks); i++ {
			artists := []string{}
			for j := 0; j < len(tracks[i].Artists); j++ {
				artists = append(artists, tracks[i].Artists[j].Name)
			}

			outputTracks = append(outputTracks, &models.SptfyTrack{
				Name:    tracks[i].Name,
				Id:      tracks[i].Id,
				Artists: artists,
				Album:   tracks[i].Album.Name,
			})
		}
		for t := 0; t < len(outputTracks); t++ {
			fmt.Printf("%s - %s - %s\n", outputTracks[t].Name, outputTracks[t].Artists, outputTracks[t].Album)
		}
		// TODO add support for pagination
	default:
		log.Fatal("argument not defined")
	}
}
