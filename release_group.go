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

// ReleaseGroup groups several different releases into a single logical entity.
// Every release belongs to one, and only one release group. More informations
// at https://musicbrainz.org/doc/Release_Group
type ReleaseGroup struct {
	ID           MBID         `xml:"id,attr"`
	Type         string       `xml:"type,attr"`
	PrimaryType  string       `xml:"primary-type"`
	Title        string       `xml:"title"`
	ArtistCredit ArtistCredit `xml:"artist-credit"`
	Releases     []*Release   `xml:"release-list>release"` // FIXME if important unmarshal count,attr
	Tags         []*Tag       `xml:"tag-list>tag"`
}

func (mbe *ReleaseGroup) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name      `xml:"metadata"`
		Ptr     *ReleaseGroup `xml:"release-group"`
	}
	res.Ptr = mbe
	return &res
}

func (mbe *ReleaseGroup) apiEndpoint() string {
	return "/release-group"
}

func (mbe *ReleaseGroup) Id() MBID {
	return mbe.ID
}

// LookupReleaseGroup performs a release-group lookup request for the given MBID.
func (c *WS2Client) LookupReleaseGroup(id MBID, inc ...string) (*ReleaseGroup, error) {
	a := &ReleaseGroup{ID: id}
	err := c.Lookup(a, inc...)

	return a, err
}

// SearchReleaseGroup queries MusicBrainz´ Search Server for ReleaseGroups.
//
// Possible search fields to provide in searchTerm are:
//
//	arid                MBID of the release group’s artist
//	artist              release group artist as it appears on the cover (Artist Credit)
//	artistname          “real name” of any artist that is included in the release group’s artist credit
//	comment             release group comment to differentiate similar release groups
//	creditname          name of any artist in multi-artist credits, as it appears on the cover.
//	primarytype         primary type of the release group (album, single, ep, other)
//	rgid                MBID of the release group
//	releasegroup        name of the release group
//	releasegroupaccent  name of the releasegroup with any accent characters retained
//	releases            number of releases in this release group
//	release             name of a release that appears in the release group
//	reid                MBID of a release that appears in the release group
//	secondarytype       secondary type of the release group (audiobook, compilation, interview, live, remix soundtrack, spokenword)
//	status              status of a release that appears within the release group
//	tag                 a tag that appears on the release group
//	type                type of the release group, old type mapping for when we did not have separate primary and secondary types
//
// With no fields specified searchTerm searches the releasgroup field only. For
// more information visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Release_Group
func (c *WS2Client) SearchReleaseGroup(searchTerm string, limit, offset int) (*ReleaseGroupSearchResponse, error) {

	result := releaseGroupListResult{}
	err := c.searchRequest("/release-group", &result, searchTerm, limit, offset)

	rsp := ReleaseGroupSearchResponse{}
	rsp.WS2ListResponse = result.ReleaseGroupList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.ReleaseGroupList.ReleaseGroups {
		rsp.ReleaseGroups = append(rsp.ReleaseGroups, v.ReleaseGroup)
		rsp.Scores[rsp.ReleaseGroups[i]] = v.Score
	}

	return &rsp, err
}

// ReleaseGroupSearchResponse is the response type returned by release group request
// methods.
type ReleaseGroupSearchResponse struct {
	WS2ListResponse
	ReleaseGroups []*ReleaseGroup
	Scores        ScoreMap
}

// ResultsWithScore returns a slice of ReleaseGroups with a specific score.
func (r *ReleaseGroupSearchResponse) ResultsWithScore(score int) []*ReleaseGroup {
	var res []*ReleaseGroup
	for _, v := range r.ReleaseGroups {
		if r.Scores[v] == score {
			res = append(res, v)
		}
	}
	return res
}

type releaseGroupListResult struct {
	ReleaseGroupList struct {
		WS2ListResponse
		ReleaseGroups []struct {
			*ReleaseGroup
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"release-group"`
	} `xml:"release-group-list"`
}
