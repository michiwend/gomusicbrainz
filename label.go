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

// Label represents an imprint, a record company or a music group. Labels refer
// mainly to imprints in MusicBrainz. Visit https://musicbrainz.org/doc/Label
// for more informations.
type Label struct {
	ID             MBID     `xml:"id,attr"`
	Name           string   `xml:"name"`
	Type           string   `xml:"type,attr"`
	SortName       string   `xml:"sort-name"`
	Disambiguation string   `xml:"disambiguation"`
	CountryCode    string   `xml:"country"`
	Area           Area     `xml:"area"`
	LabelCode      int      `xml:"label-code"`
	Lifespan       Lifespan `xml:"life-span"`
	Aliases        []*Alias `xml:"alias-list>alias"`
}

// LabelSearchResponse is the response type returned by the label search method.
type LabelSearchResponse struct {
	WS2ListResponse
	Labels []*Label
	Scores ScoreMap
}

// ResultsWithScore returns a slice of Labels with a specific score.
func (r *LabelSearchResponse) ResultsWithScore(score int) []*Label {
	var res []*Label
	for k, v := range r.Scores {
		if v == score {
			res = append(res, k.(*Label))
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
