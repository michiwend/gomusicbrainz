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

// CDStub represents an anonymously submitted track list.
type CDStub struct {
	ID        string `xml:"id,attr"` // seems not to be a valid MBID (UUID)
	Title     string `xml:"title"`
	Artist    string `xml:"artist"`
	Barcode   string `xml:"barcode"`
	Comment   string `xml:"comment"`
	TrackList struct {
		Count int `xml:"count,attr"`
	} `xml:"track-list"`
}

// CDStubSearchResponse is the response type returned by the CDStub search method.
type CDStubSearchResponse struct {
	WS2ListResponse
	CDStubs []*CDStub
	Scores  ScoreMap
}

// ResultsWithScore returns a slice of CDStubs with a specific score.
func (r *CDStubSearchResponse) ResultsWithScore(score int) []*CDStub {
	var res []*CDStub
	for k, v := range r.Scores {
		if v == score {
			res = append(res, k.(*CDStub))
		}
	}
	return res
}

type cdStubListResult struct {
	CDStubList struct {
		WS2ListResponse
		CDStubs []struct {
			*CDStub
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"cdstub"`
	} `xml:"cdstub-list"`
}
