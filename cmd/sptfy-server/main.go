package main

import (
	"github.com/llater/sptfy/src/sptfy"
	"log"
	"os"
)

const (
	CONFIGURATION_FILE_PATH   = "~/.sptfy/config"
	CLIENT_ID_ENVVAR_NAME     = "SPOTIFY_CLIENT_ID"
	CLIENT_SECRET_ENVVAR_NAME = "SPOTIFY_CLIENT_SECRET"
)

type sptfyConfig struct {
	spotifyClientId     string
	spotifyClientSecret string
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

	server, err := sptfy.NewSptfyServer(clientId, clientSecret)
	if err != nil {
		log.Println("g")
		log.Fatal(err)
	}

	// Verify connection
	err := server.Up()

	if err != nil {
		log.Println("sptfy server is not up")
		log.Fatal(err)
	}
}
