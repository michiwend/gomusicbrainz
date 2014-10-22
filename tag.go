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

// Tag is the common type for Tags.
type Tag struct {
	Count int    `xml:"count,attr"`
	Name  string `xml:"name"`
}

// SearchTag queries MusicBrainz' Search Server for Tags.
// searchTerm only contains the tag field. For more information visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Tag
func (c *WS2Client) SearchTag(searchTerm string, limit, offset int) (*TagSearchResponse, error) {

	result := tagListResult{}
	err := c.searchRequest("/tag", &result, searchTerm, limit, offset)

	rsp := TagSearchResponse{}
	rsp.WS2ListResponse = result.TagList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.TagList.Tags {
		rsp.Tags = append(rsp.Tags, v.Tag)
		rsp.Scores[rsp.Tags[i]] = v.Score
	}

	return &rsp, err
}

// TagSearchResponse is the response type returned by the SearchTag method.
type TagSearchResponse struct {
	WS2ListResponse
	Tags   []*Tag
	Scores ScoreMap
}

// ResultsWithScore returns a slice of Tags with a specific score.
func (r *TagSearchResponse) ResultsWithScore(score int) []*Tag {
	var res []*Tag
	for k, v := range r.Scores {
		if v == score {
			res = append(res, k.(*Tag))
		}
	}
	return res
}

type tagListResult struct {
	TagList struct {
		WS2ListResponse
		Tags []struct {
			*Tag
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"tag"`
	} `xml:"tag-list"`
}
