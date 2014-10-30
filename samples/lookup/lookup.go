package main

import (
	"fmt"

	"github.com/michiwend/gomusicbrainz"
)

func main() {

	// create a new WS2Client.
	client, _ := gomusicbrainz.NewWS2Client(
		"https://musicbrainz.org/ws/2",
		"A GoMusicBrainz example",
		"0.0.1-beta",
		"http://github.com/michiwend/gomusicbrainz")

	// Lookup artist by id.
	artist, err := client.LookupArtist("9a709693-b4f8-4da9-8cc1-038c911a61be")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v", artist)
}
