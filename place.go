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

// PlaceSearchResponse is the response type returned by the place search method.
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
