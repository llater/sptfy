package main

import (
	"errors"
	"github.com/llater/sptfy/pkg/clients"
	"log"
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
		log.Fatal("Must supply an argument. --help for help")
	} else {
		command = os.Args[1] // Call only the first argument
	}

	switch command {
	case "me":
		me, err := client.Me()
		crash(err)
		log.Printf("Name: %s\nEmail: %s\nSpotifyID: %s", me.DisplayName, me.Email, me.Id)
	case "search":
		log.Print("work on promptui")
	default:
		log.Fatal("argument not defined")
	}
}
