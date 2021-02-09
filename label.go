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

// LabelInfo contains a label and links it to a catalog number.
type LabelInfo struct {
	CatalogNumber string `xml:"catalog-number"`
	Label         *Label `xml:"label"`
}

// Label represents an imprint, a record company or a music group. Labels refer
// mainly to imprints in MusicBrainz. Visit https://musicbrainz.org/doc/Label
// for more information.
type Label struct {
	ID             MBID     `xml:"id,attr"`
	Name           string   `xml:"name"`
	Type           string   `xml:"type,attr"`
	TypeId         string   `xml:"type-id,attr"`
	SortName       string   `xml:"sort-name"`
	Disambiguation string   `xml:"disambiguation"`
	CountryCode    string   `xml:"country"`
	Area           Area     `xml:"area"`
	LabelCode      int      `xml:"label-code"`
	Lifespan       Lifespan `xml:"life-span"`
	Aliases        []*Alias `xml:"alias-list>alias"`
}

func (mbe *Label) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name `xml:"metadata"`
		Ptr     *Label   `xml:"label"`
	}
	res.Ptr = mbe
	return &res
}

func (mbe *Label) apiEndpoint() string {
	return "/label"
}

func (mbe *Label) Id() MBID {
	return mbe.ID
}

// LookupLabel performs a label lookup request for the given MBID.
func (c *WS2Client) LookupLabel(id MBID, inc ...string) (*Label, error) {
	a := &Label{ID: id}
	err := c.Lookup(a, inc...)

	return a, err
}

// SearchLabel queries MusicBrainz´ Search Server for Labels.
//
// Possible search fields to provide in searchTerm are:
//
//	alias        the aliases/misspellings for this label
//	area         label area
//	begin        label founding date
//	code         label code (only the figures part, i.e. without "LC")
//	comment      label comment to differentiate similar labels
//	country      The two letter country code of the label country
//	end          label dissolution date
//	ended        true if know ended even if do not know end date
//	ipi          ipi
//	label        label name
//	labelaccent  name of the label with any accent characters retained
//	laid         MBID of the label
//	sortname     label sortname
//	type         label type
//	tag          folksonomy tag
//
// With no fields specified searchTerm searches the label, sortname and alias
// fields. For more information visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Label
func (c *WS2Client) SearchLabel(searchTerm string, limit, offset int) (*LabelSearchResponse, error) {

	result := labelListResult{}
	err := c.searchRequest("/label", &result, searchTerm, limit, offset)

	rsp := LabelSearchResponse{}
	rsp.WS2ListResponse = result.LabelList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.LabelList.Labels {
		rsp.Labels = append(rsp.Labels, v.Label)
		rsp.Scores[rsp.Labels[i]] = v.Score
	}

	return &rsp, err
}

// LabelSearchResponse is the response type returned by the SearchLabel method.
type LabelSearchResponse struct {
	WS2ListResponse
	Labels []*Label
	Scores ScoreMap
}

// ResultsWithScore returns a slice of Labels with a min score.
func (r *LabelSearchResponse) ResultsWithScore(score int) []*Label {
	var res []*Label
	for _, v := range r.Labels {
		if r.Scores[v] >= score {
			res = append(res, v)
		}
	}
	return res
}

type labelListResult struct {
	LabelList struct {
		WS2ListResponse
		Labels []struct {
			*Label
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"label"`
	} `xml:"label-list"`
}
