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

// ReleaseGroup groups several different releases into a single logical entity.
// Every release belongs to one, and only one release group. More informations
// at https://musicbrainz.org/doc/Release_Group
type ReleaseGroup struct {
	ID           string       `xml:"id,attr"`
	Type         string       `xml:"type,attr"`
	PrimaryType  string       `xml:"primary-type"`
	Title        string       `xml:"title"`
	ArtistCredit ArtistCredit `xml:"artist-credit"`
	Releases     []Release    `xml:"release-list>release"` // FIXME if important unmarshal count,attr
	Tags         []Tag        `xml:"tag-list>tag"`
}

// ReleaseGroupResponse is the response type returned by release group request
// methods.
type ReleaseGroupResponse struct {
	WS2ListResponse
	ReleaseGroups []ReleaseGroup
	Scores        ScoreMap
}

// ResultsWithScore returns a slice of ReleaseGroups with a specific score.
func (r *ReleaseGroupResponse) ResultsWithScore(score int) []ReleaseGroup {
	var res []ReleaseGroup
	for k, v := range r.Scores {
		if v == score {
			res = append(res, *k.(*ReleaseGroup))
		}
	}
	return res
}

type releaseGroupListResult struct {
	ReleaseGroupList struct {
		WS2ListResponse
		ReleaseGroups []struct {
			ReleaseGroup
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"release-group"`
	} `xml:"release-group-list"`
}
