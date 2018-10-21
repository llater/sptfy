package main

import (
	"github.com/llater/sptfy/src/sptfy"
	"log"
	"os"
)

const (
	CONFIGURATION_FILE_PATH   = "~/.sptfy/config"
	REFRESH_TOKEN_FILE_PATH   = "~/.sptfy/token"
	CLIENT_ID_ENVVAR_NAME     = "SPOTIFY_CLIENT_ID"
	CLIENT_SECRET_ENVVAR_NAME = "SPOTIFY_CLIENT_SECRET"
	REDIRECT_URI_ENVVAR_NAME  = "OAUTH2_CLIENT_REDIRECT_URI"
	REFRESH_TOKEN_ENVVAR_NAME = "SPOTIFY_REFRESH_TOKEN"
)

var ()

type sptfyConfig struct {
	spotifyClientId string
}

func readConfigFromFile(configPath string) (*sptfyConfig, error) {
	base := sptfyConfig{}
	return &base, nil
}

func main() {

	// Read environment variables
	clientId, ok := os.LookupEnv(CLIENT_ID_ENVVAR_NAME)
	if !ok {
		log.Fatal("failed to read OAuth2 client id from the environment")
	}

	clientSecret, ok := os.LookupEnv(CLIENT_SECRET_ENVVAR_NAME)
	if !ok {
		log.Fatal("failed to read OAuth2 client secret from the environment")
	}

	redirectUri, ok := os.LookupEnv(REDIRECT_URI_ENVVAR_NAME)
	if !ok {
		log.Fatal("failed to read redirect URI from the environment")
	}

	spotify, err := sptfy.NewSpotifyApiClient(clientId, clientSecret, redirectUri)
	if err != nil {
		log.Println("g")
		log.Fatal(err)
	}

	// Verify connection
	success, err := spotify.Ping()

	if err != nil {
		log.Println("Failed to ping spotify")
		log.Fatal(err)
	}

	log.Printf("Logged in as %s", success)

	/*	var command string
		if len(os.Args) < 2 {
			log.Fatal("Must supply an argument. --help for help")
		} else {
			command = os.Args[1] // Call only the first argument
		}

		switch command {
		case "me": */
	/*default:
		log.Fatal("argument not defined")
	}
	*/
	log.Println("finished executing")
}
