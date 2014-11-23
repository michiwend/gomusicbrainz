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

type Recording struct {
	ID             MBID         `xml:"id,attr"`
	Title          string       `xml:"title"`
	Length         int          `xml:"length"`
	Disambiguation string       `xml:"disambiguation"`
	ArtistCredit   ArtistCredit `xml:"artist-credit"`

	// TODO add refs
}

func (mbe *Recording) lookupResult() interface{} {
	var res struct {
		XMLName xml.Name   `xml:"metadata"`
		Ptr     *Recording `xml:"recording"`
	}
	res.Ptr = mbe
	return &res
}

func (mbe *Recording) apiEndpoint() string {
	return "/recording"
}

func (mbe *Recording) Id() MBID {
	return mbe.ID
}

// LookupRecording performs an recording lookup request for the given MBID.
func (c *WS2Client) LookupRecording(id MBID, inc ...string) (*Recording, error) {
	a := &Recording{ID: id}
	err := c.Lookup(a, inc...)

	return a, err
}

// SearchRecording queries MusicBrainz´ Search Server for Recordings.
//
// Possible search fields to provide in searchTerm are:
//
//	arid             artist id
//	artist           artist name is name(s) as it appears on the recording
//	artistname       an artist on the recording, each artist added as a separate field
//	creditname       name credit on the recording, each artist added as a separate field
//	comment          recording disambiguation comment
//	country          recording release country
//	date             recording release date
//	dur              duration of track in milliseconds
//	format           recording release format
//	isrc             ISRC of recording
//	number           free text track number
//	position         the medium that the recording should be found on, first medium is position 1
//	primarytype      primary type of the release group (album, single, ep, other)
//	puid             PUID of recording
//	qdur             quantized duration (duration / 2000)
//	recording        name of recording or a track associated with the recording
//	recordingaccent  name of the recording with any accent characters retained
//	reid             release id
//	release          release name
//	rgid             release group id
//	rid              recording id
//	secondarytype    secondary type of the release group (audiobook, compilation, interview, live, remix soundtrack, spokenword)
//	status           Release status (official, promotion, Bootleg, Pseudo-Release)
//	tid              track id
//	tnum             track number on medium
//	tracks           number of tracks in the medium on release
//	tracksrelease    number of tracks on release as a whole
//	tag              folksonomy tag
//	type             type of the release group, old type mapping for when we did not have separate primary and secondary types or use standalone for standalone recordings
//	video            true to only show video tracks
//
// With no fields specified searchTerm searches the recording field only. For
// more information visit
// http://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Recording
func (c *WS2Client) SearchRecording(searchTerm string, limit, offset int) (*RecordingSearchResponse, error) {

	result := recordingListResult{}
	err := c.searchRequest("/recording", &result, searchTerm, limit, offset)

	rsp := RecordingSearchResponse{}
	rsp.WS2ListResponse = result.RecordingList.WS2ListResponse
	rsp.Scores = make(ScoreMap)

	for i, v := range result.RecordingList.Recordings {
		rsp.Recordings = append(rsp.Recordings, v.Recording)
		rsp.Scores[rsp.Recordings[i]] = v.Score
	}

	return &rsp, err
}

// RecordingSearchResponse is the response type returned by the SearchRecording
// method.
type RecordingSearchResponse struct {
	WS2ListResponse
	Recordings []*Recording
	Scores     ScoreMap
}

// ResultsWithScore returns a slice of Recordings with a specific score.
func (r *RecordingSearchResponse) ResultsWithScore(score int) []*Recording {
	var res []*Recording
	for k, v := range r.Scores {
		if v == score {
			res = append(res, k.(*Recording))
		}
	}
	return res
}

type recordingListResult struct {
	RecordingList struct {
		WS2ListResponse
		Recordings []struct {
			*Recording
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"recording"`
	} `xml:"recording-list"`
}
