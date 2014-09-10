package gomusicbrainz

import (
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
)

func New(apiBaseURL string) *GoMusicBrainz {

	c := GoMusicBrainz{
		APIBase: apiBaseURL,
	}

	return &c
}

type GoMusicBrainz struct {
	APIBase string
}

func (c *GoMusicBrainz) getReqeust(params url.Values, endpoint string) {

}

func (c *GoMusicBrainz) SearchArtist(query string) ([]Artist, error) {

	endpoint := "artist/"

	params := url.Values{
		"query": {query},
	}

	resp, err := http.Get(c.APIBase + endpoint + "?" + params.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)

	var result ArtistSearchRequest
	if err = decoder.Decode(&result); err != nil {
		return []Artist{}, err
	}

	return result.ArtistList.Artists, nil
}
