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

// Annotation is a miniature wiki that can be added to any existing artists,
// labels, recordings, releases, release groups and works. More informations at
// https://musicbrainz.org/doc/Annotation
type Annotation struct {
	Type   string `xml:"type,attr"`
	Entity string `xml:"entity"`
	Name   string `xml:"name"`
	Text   string `xml:"text"`
}

// SearchAnnotation queries MusicBrainz´ Search Server for Annotations.
//
// Possible search fields to provide in searchTerm are:
//
//	text    The content of the annotation
//	type    The entity type (artist, releasegroup, release, recording, work, label)
//	name    The name of the entity
//	entity  The entity's MBID
//
// For more information visit
// http://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Annotation
func (c *WS2Client) SearchAnnotation(searchTerm string, limit, offset int) (*AnnotationSearchResponse, error) {

	result := annotationListResult{}
	err := c.searchRequest("/annotation", &result, searchTerm, limit, offset)

	rsp := AnnotationSearchResponse{}
	rsp.WS2ListResponse = result.AnnotationList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.AnnotationList.Annotations {
		rsp.Annotations = append(rsp.Annotations, v.Annotation)
		rsp.Scores[rsp.Annotations[i]] = v.Score
	}

	return &rsp, err
}

// AnnotationSearchResponse is the response type returned by annotation request
// methods.
type AnnotationSearchResponse struct {
	WS2ListResponse
	Annotations []*Annotation
	Scores      ScoreMap
}

// ResultsWithScore returns a slice of Annotations with a specific score.
func (r *AnnotationSearchResponse) ResultsWithScore(score int) []*Annotation {
	var res []*Annotation
	for k, v := range r.Scores {
		if v == score {
			res = append(res, k.(*Annotation))
		}
	}
	return res
}

type annotationListResult struct {
	AnnotationList struct {
		WS2ListResponse
		Annotations []struct {
			*Annotation
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"annotation"`
	} `xml:"annotation-list"`
}
