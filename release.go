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

// Release represents a unique release (i.e. issuing) of a product on a
// specific date with specific release information such as the country, label,
// barcode, packaging, etc. More information at https://musicbrainz.org/doc/Release
type Release struct {
	ID                 MBID               `xml:"id,attr"`
	Title              string             `xml:"title"`
	Status             string             `xml:"status"`
	Disambiguation     string             `xml:"disambiguation"`
	TextRepresentation TextRepresentation `xml:"text-representation"`
	ArtistCredit       ArtistCredit       `xml:"artist-credit"`
	ReleaseGroup       ReleaseGroup       `xml:"release-group"`
	Date               BrainzTime         `xml:"date"`
	CountryCode        string             `xml:"country"`
	Barcode            string             `xml:"barcode"`
	Asin               string             `xml:"asin"`
	Quality            string             `xml:"quality"`
	LabelInfos         []LabelInfo        `xml:"label-info-list>label-info"`
	Mediums            []*Medium          `xml:"medium-list>medium"`
	Relations          TargetRelationsMap `xml:"relation-list"`
}

func (mbe *Release) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name `xml:"metadata"`
		Ptr     *Release `xml:"release"`
	}
	res.Ptr = mbe
	return &res
}

func (mbe *Release) apiEndpoint() string {
	return "/release"
}

func (mbe *Release) Id() MBID {
	return mbe.ID
}

// LookupRelease performs a release lookup request for the given MBID.
func (c *WS2Client) LookupRelease(id MBID, inc ...string) (*Release, error) {
	a := &Release{ID: id}
	err := c.Lookup(a, inc...)

	return a, err
}

// SearchRelease queries MusicBrainz´ Search Server for Releases.
// With no fields specified searchTerm searches the release field only. For a
// list of all valid fields visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Release
func (c *WS2Client) SearchRelease(searchTerm string, limit, offset int) (*ReleaseSearchResponse, error) {

	result := releaseListResult{}
	err := c.searchRequest("/release", &result, searchTerm, limit, offset)

	rsp := ReleaseSearchResponse{}
	rsp.WS2ListResponse = result.ReleaseList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.ReleaseList.Releases {
		rsp.Releases = append(rsp.Releases, v.Release)
		rsp.Scores[rsp.Releases[i]] = v.Score
	}

	return &rsp, err
}

// ReleaseSearchResponse is the response type returned by the SearchRelease method.
type ReleaseSearchResponse struct {
	WS2ListResponse
	Releases []*Release
	Scores   ScoreMap
}

// ResultsWithScore returns a slice of Releases with a specific score.
func (r *ReleaseSearchResponse) ResultsWithScore(score int) []*Release {
	var res []*Release
	for k, v := range r.Scores {
		if v == score {
			res = append(res, k.(*Release))
		}
	}
	return res
}

type releaseListResult struct {
	ReleaseList struct {
		WS2ListResponse
		Releases []struct {
			*Release
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"release"`
	} `xml:"release-list"`
}
