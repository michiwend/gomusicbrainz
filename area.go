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

// Area represents a geographic region or settlement.
type Area struct {
	ID            MBID           `xml:"id,attr"`
	Type          string         `xml:"type,attr"`
	Name          string         `xml:"name"`
	SortName      string         `xml:"sort-name"`
	ISO31661Codes []ISO31661Code `xml:"iso-3166-1-code-list>iso-3166-1-code"`
	Lifespan      Lifespan       `xml:"life-span"`
	Aliases       []Alias        `xml:"alias-list>alias"`
}

type ISO31661Code string

func (mbe *Area) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name `xml:"metadata"`
		Ptr     *Area    `xml:"area"`
	}
	res.Ptr = mbe
	return &res
}

func (mbe *Area) apiEndpoint() string {
	return "/area"
}

func (mbe *Area) Id() MBID {
	return mbe.ID
}

// LookupArea performs an area lookup request for the given MBID.
func (c *WS2Client) LookupArea(id MBID, inc ...string) (*Area, error) {
	a := &Area{ID: id}
	err := c.Lookup(a, inc...)

	return a, err
}

// SearchArea queries MusicBrainz´ Search Server for Areas.
//
// Possible search fields to provide in searchTerm are:
//
//	aid       the area ID
//	alias     the aliases/misspellings for this area
//	area      area name
//	begin     area begin date
//	comment   disambugation comment
//	end 	  area end date
//	ended 	  area ended
//	sortname  area sort name
//	iso       area iso1, iso2 or iso3 codes
//	iso1      area iso1 codes
//	iso2      area iso3 codes
//	iso3      area iso3 codes
//	type      the aliases/misspellings for this label
//
// With no fields specified searchTerm searches the area and sortname fields.
// For more information visit
// http://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Area
func (c *WS2Client) SearchArea(searchTerm string, limit, offset int) (*AreaSearchResponse, error) {

	result := areaListResult{}
	err := c.searchRequest("/area", &result, searchTerm, limit, offset)

	rsp := AreaSearchResponse{}
	rsp.WS2ListResponse = result.AreaList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.AreaList.Areas {
		rsp.Areas = append(rsp.Areas, v.Area)
		rsp.Scores[rsp.Areas[i]] = v.Score
	}

	return &rsp, err
}

// AreaSearchResponse is the response type returned by the SearchArea method.
type AreaSearchResponse struct {
	WS2ListResponse
	Areas  []*Area
	Scores ScoreMap
}

// ResultsWithScore returns a slice of Areas with a min score.
func (r *AreaSearchResponse) ResultsWithScore(score int) []*Area {
	var res []*Area
	for _, v := range r.Areas {
		if r.Scores[v] >= score {
			res = append(res, v)
		}
	}
	return res
}

type areaListResult struct {
	AreaList struct {
		WS2ListResponse
		Areas []struct {
			*Area
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"area"`
	} `xml:"area-list"`
}
