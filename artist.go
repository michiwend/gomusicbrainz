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

package gomusicbrainz

import "encoding/xml"

// Artist represents generally a musician, a group of musicians, a collaboration
// of multiple musicians or other music professionals.
type Artist struct {
	ID             MBID     `xml:"id,attr"`
	Type           string   `xml:"type,attr"`
	Name           string   `xml:"name"`
	Disambiguation string   `xml:"disambiguation"`
	SortName       string   `xml:"sort-name"`
	CountryCode    string   `xml:"country"`
	Lifespan       Lifespan `xml:"life-span"`
	Aliases        []*Alias `xml:"alias-list>alias"`
}

func (mble *Artist) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name `xml:"metadata"`
		Ptr     *Artist  `xml:"artist"`
	}
	res.Ptr = mble
	return &res
}

func (mble *Artist) apiEndpoint() string {
	return "/artist"
}

func (mble *Artist) id() MBID {
	return mble.ID
}

// LookupArtist performs an artist lookup request for the given MBID.
func (c *WS2Client) LookupArtist(id MBID) (*Artist, error) {
	a := &Artist{ID: id}
	err := c.Lookup(a)

	return a, err
}

// SearchArtist queries MusicBrainz´ Search Server for Artists.
// With no fields specified searchTerm searches the artist, sortname and alias
// fields. For a list of all valid fields visit
// http://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Artist
func (c *WS2Client) SearchArtist(searchTerm string, limit, offset int) (*ArtistSearchResponse, error) {

	result := artistListResult{}
	err := c.searchRequest("/artist", &result, searchTerm, limit, offset)

	rsp := ArtistSearchResponse{}
	rsp.WS2ListResponse = result.ArtistList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.ArtistList.Artists {
		rsp.Artists = append(rsp.Artists, v.Artist)
		rsp.Scores[rsp.Artists[i]] = v.Score
	}

	return &rsp, err
}

// ArtistSearchResponse is the response type returned by the SearchArtist method.
type ArtistSearchResponse struct {
	WS2ListResponse
	Artists []*Artist
	Scores  ScoreMap
}

// ResultsWithScore returns a slice of Artists with a specific score.
func (r *ArtistSearchResponse) ResultsWithScore(score int) []*Artist {
	var res []*Artist
	for k, v := range r.Scores {
		if v == score {
			res = append(res, k.(*Artist))
		}
	}
	return res
}

type artistListResult struct {
	ArtistList struct {
		WS2ListResponse
		Artists []struct {
			*Artist
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"artist"`
	} `xml:"artist-list"`
}
