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

// Place represents a building or outdoor area used for performing or producing
// music.
type Place struct {
	ID          MBID          `xml:"id,attr"`
	Type        string        `xml:"type,attr"`
	Name        string        `xml:"name"`
	Address     string        `xml:"address"`
	Coordinates MBCoordinates `xml:"coordinates"`
	Area        Area          `xml:"area"`
	Lifespan    Lifespan      `xml:"life-span"`
	Aliases     []*Alias      `xml:"alias-list>alias"`
}

func (mbe *Place) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name `xml:"metadata"`
		Ptr     *Place   `xml:"place"`
	}
	res.Ptr = mbe
	return &res
}

func (mbe *Place) apiEndpoint() string {
	return "/place"
}

func (mbe *Place) Id() MBID {
	return mbe.ID
}

// LookupPlace performs a place lookup request for the given MBID.
func (c *WS2Client) LookupPlace(id MBID, inc ...string) (*Place, error) {
	a := &Place{ID: id}
	err := c.Lookup(a, inc...)

	return a, err
}

// SearchPlace queries MusicBrainz´ Search Server for Places.
//
// Possible search fields to provide in searchTerm are:
//
//	pid       the place ID
//	address   the address of this place
//	alias     the aliases/misspellings for this area
//	area      area name
//	begin     place begin date
//	comment   disambiguation comment
//	end       place end date
//	ended     place ended
//	lat       place latitude
//	long      place longitude
//	sortname  place sort name
//	type      the aliases/misspellings for this place
//
// With no fields specified searchTerm searches the place, alias, address and
// area fields. For more information visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Place
func (c *WS2Client) SearchPlace(searchTerm string, limit, offset int) (*PlaceSearchResponse, error) {

	result := placeListResult{}
	err := c.searchRequest("/place", &result, searchTerm, limit, offset)

	rsp := PlaceSearchResponse{}
	rsp.WS2ListResponse = result.PlaceList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.PlaceList.Places {
		rsp.Places = append(rsp.Places, v.Place)
		rsp.Scores[rsp.Places[i]] = v.Score
	}

	return &rsp, err
}

// PlaceSearchResponse is the response type returned by the SearchPlace method.
type PlaceSearchResponse struct {
	WS2ListResponse
	Places []*Place
	Scores ScoreMap
}

// ResultsWithScore returns a slice of Places with a specific score.
func (r *PlaceSearchResponse) ResultsWithScore(score int) []*Place {
	var res []*Place
	for k, v := range r.Scores {
		if v == score {
			res = append(res, k.(*Place))
		}
	}
	return res
}

type placeListResult struct {
	PlaceList struct {
		WS2ListResponse
		Places []struct {
			*Place
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"place"`
	} `xml:"place-list"`
}
