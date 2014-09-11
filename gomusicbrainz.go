/*
 *	Copyright (c) 2014 Michael Wendland
 *
 *	Permission is hereby granted, free of charge, to any person obtaining a
 *	copy of this software and associated documentation files (the "Software"),
 *	to deal in the Software without restriction, including without limitation
 *	the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *	and/or sell copies of the Software, and to permit persons to whom the
 *	Software is furnished to do so, subject to the following conditions:
 *
 *	The above copyright notice and this permission notice shall be included in
 *	all copies or substantial portions of the Software.
 *
 *	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 *	FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 *	IN THE SOFTWARE.
 *
 *	Authors:
 *		Michael Wendland <michael@michiwend.com>
 */

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
