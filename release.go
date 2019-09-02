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

import (
	"encoding/xml"
	"net/url"
	"path"
)

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

type DiscIdRelease struct {
	Release
}

func (mbe *Release) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name `xml:"metadata"`
		Ptr     *Release `xml:"release"`
	}
	res.Ptr = mbe
	return &res
}

type LookupDiscIdResponse struct {
	XMLName xml.Name `xml:"metadata"`
	id MBID
	Releases []*Release `xml:"release-list>release"`
}

func (l LookupDiscIdResponse) Id() MBID {
	return l.id
}

func (LookupDiscIdResponse) apiEndpoint() string {
	return "discid"
}

func (d *LookupDiscIdResponse) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name `xml:"metadata"`
		Ptr *[]*Release `xml:"release-list>release"`
	}
	res.Ptr = &d.Releases
	return &res
}

func (c *WS2Client) LookupDiscId(toc string, inc ...string) ([]*Release, error) {
	result := &LookupDiscIdResponse{id:"-"}
	err := c.getRequest(result.lookupResult(), url.Values{"toc": []string{toc}, "cdstubs": []string{"no"}},
		path.Join(
			result.apiEndpoint(),
			string(result.Id()),
		),
	)
	return result.Releases, err
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
//
// Possible search fields to provide in searchTerm are:
//
//	arid           artist id
//	artist         complete artist name(s) as it appears on the release
//	artistname     an artist on the release, each artist added as a separate field
//	asin           the Amazon ASIN for this release
//	barcode        The barcode of this release
//	catno          The catalog number for this release, can have multiples when major using an imprint
//	comment        Disambiguation comment
//	country        The two letter country code for the release country
//	creditname     name credit on the release, each artist added as a separate field
//	date           The release date (format: YYYY-MM-DD)
//	discids        total number of cd ids over all mediums for the release
//	discidsmedium  number of cd ids for the release on a medium in the release
//	format         release format
//	laid           The label id for this release, a release can have multiples when major using an imprint
//	label          The name of the label for this release, can have multiples when major using an imprint
//	lang           The language for this release. Use the three character ISO 639 codes to search for a specific language. (e.g. lang:eng)
//	mediums        number of mediums in the release
//	primarytype    primary type of the release group (album, single, ep, other)
//	puid           The release contains recordings with these puids
//	quality        The quality of the release (low, normal, high)
//	reid           release id
//	release        release name
//	releaseaccent  name of the release with any accent characters retained
//	rgid           release group id
//	script         The 4 character script code (e.g. latn) used for this release
//	secondarytype  secondary type of the release group (audiobook, compilation, interview, live, remix, soundtrack, spokenword)
//	status         release status (e.g official)
//	tag            a tag that appears on the release
//	tracks         total number of tracks over all mediums on the release
//	tracksmedium   number of tracks on a medium in the release
//	type           type of the release group, old type mapping for when we did not have separate primary and secondary types
//
// With no fields specified searchTerm searches the release field only. For
// more information visit
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

// ResultsWithScore returns a slice of Releases with a min score.
func (r *ReleaseSearchResponse) ResultsWithScore(score int) []*Release {
	var res []*Release
	for _, v := range r.Releases {
		if r.Scores[v] >= score {
			res = append(res, v)
		}
	}
	return res
}

// OriginalRelease is a helper function that returns the earliest release of
// a release array with the most accurate date. It can be used to determine
// the original/first release from releases of a release group.
func OriginalRelease(releases []*Release) *Release {

	if len(releases) == 0 {
		return nil
	}
	original := releases[0] // fall back on the first item

	for _, release := range releases {

		if !release.Date.IsZero() {

			if release.Date.Year() < original.Date.Year() || original.Date.IsZero() {
				original = release
			} else if release.Date.Year() == original.Date.Year() &&
				release.Date.Accuracy > Year {

				if original.Date.Accuracy == Year ||
					release.Date.Month() < original.Date.Month() {

					original = release

				} else if release.Date.Month() == original.Date.Month() &&
					release.Date.Accuracy > Month {

					if original.Date.Accuracy == Month ||
						release.Date.Day() < original.Date.Day() {
						original = release
					}
				}
			}
		}
	}

	return original
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
