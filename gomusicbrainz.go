/*
 * Copyright (c) 2014 Michael Wendland
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 * IN THE SOFTWARE.
 *
 * 	Authors:
 * 		Michael Wendland <michael@michiwend.com>
 */

/*
Package gomusicbrainz implements a MusicBrainz WS2 client library.
*/
package gomusicbrainz

import (
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func NewWS2Client() *GoMusicBrainz {
	c := GoMusicBrainz{}
	c.WS2RootURL, _ = url.Parse("https://musicbrainz.org/ws/2")
	return &c
}

type GoMusicBrainz struct {
	WS2RootURL *url.URL
}

func (c *GoMusicBrainz) getReqeust(data interface{}, params url.Values, endpoint string) error {

	resp, err := http.Get(c.WS2RootURL.String() + endpoint + "?" + params.Encode())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)

	if err = decoder.Decode(data); err != nil {
		return err
	}
	return nil
}

// SearchArtist queries MusicBrainz' Search Server for Artists.
// searchTerm follows the Apache Lucene syntax. If no fields were specified the
// Search Server searches for searchTerm in any of the fields artist, sortname
// and alias. For a list of all valid search fields visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Artist
// limit defines how many entries will be returned by the server (allowed
// range 1-100, defaults to 25). offset can be used for result pagination.
func (c *GoMusicBrainz) SearchArtist(searchTerm string, limit int, offset int) ([]Artist, error) {

	result := artistSearchRequest{}
	endpoint := "/artist"

	err := c.getReqeust(&result, url.Values{
		"query":  {searchTerm},
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}, endpoint)

	if err != nil {
		return []Artist{}, err
	}

	return result.ArtistList.Artists, nil
}

// SearchRelease queries MusicBrainz' Search Server for Releases.
// searchTerm follows the Apache Lucene syntax. If no fields were specified the
// Search Server searches the release field only. For a list of all valid
// search fields visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Release
// limit defines how many entries will be returned by the server (allowed
// range 1-100, defaults to 25). offset can be used for result pagination.
func (c *GoMusicBrainz) SearchRelease(searchTerm string, limit int, offset int) ([]Release, error) {

	result := releaseSearchRequest{}
	endpoint := "/release"

	err := c.getReqeust(&result, url.Values{
		"query":  {searchTerm},
		"limit":  {strconv.Itoa(limit)},
		"offset": {strconv.Itoa(offset)},
	}, endpoint)

	if err != nil {
		return []Release{}, err
	}

	return result.ReleaseList.Releases, nil
}
