package main

import (
	"errors"
	"fmt"
	"github.com/llater/sptfy/pkg/clients"
	"github.com/llater/sptfy/pkg/models"
	"github.com/llater/sptfy/pkg/utils"
	"github.com/manifoldco/promptui"
	"log"
	"net/url"
	"os"
)

const (
	CLIENT_ID_ENVVAR_NAME     = "SPOTIFY_CLIENT_ID"
	CLIENT_SECRET_ENVVAR_NAME = "SPOTIFY_CLIENT_SECRET"
)

func crash(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	// TODO Change this to read from filepath
	// TODO Select type of client: logged-in or query
	clientId, ok := os.LookupEnv(CLIENT_ID_ENVVAR_NAME)
	if !ok {
		panic(errors.New("Spotify client ID not found"))
	}

	clientSecret, ok := os.LookupEnv(CLIENT_SECRET_ENVVAR_NAME)
	if !ok {
		panic(errors.New("Spotify client secret not found"))
	}

	client, err := clients.NewSpotifyClientCredentialsClient(clientId, clientSecret)
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
		for t := 0; len(outputTracks); t++ {
			fmt.Printf("%s - %s - %s\n", outputTracks[t].Name, outputTracks[t].Artists, outputTracks[t].Album)
		}
		// TODO add support for pagination
	default:
		log.Fatal("argument not defined")
	}
}
