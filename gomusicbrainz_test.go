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
 *	Authors:
 * 		Michael Wendland <michael@michiwend.com>
 */

package gomusicbrainz

// TODO use testdata from https://github.com/metabrainz/mmd-schema/tree/master/test-data/valid

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aryann/difflib"
	"github.com/davecgh/go-spew/spew"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *WS2Client
)

// Init multiplexer and httptest server
func setupHTTPTesting() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewWS2Client(
		server.URL,
		"Application Name",
		"Version",
		"Contact",
	)
}

// handleFunc passes response to the http client.
/*
func handleFunc(url string, response *string, t *testing.T) {
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, *response)
	})
}*/

// serveTestFile responses to the http client with content of a test file
// located in ./testdata
func serveTestFile(url string, testfile string, t *testing.T) {

	//TODO check request URL if it matches one of the following patterns
	//lookup:   /<ENTITY>/<MBID>?inc=<INC>
	//browse:   /<ENTITY>?<ENTITY>=<MBID>&limit=<LIMIT>&offset=<OFFSET>&inc=<INC>
	//search:   /<ENTITY>?query=<QUERY>&limit=<LIMIT>&offset=<OFFSET>

	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if testing.Verbose() {
			fmt.Println("GET request:", r.URL.String())
		}

		http.ServeFile(w, r, "./testdata/"+testfile)
	})
}

// pretty prints a diff
func requestDiff(want, returned interface{}) string {
	spew.Config.SortKeys = true
	spew.Config.ContinueOnMethod = true
	recs := difflib.Diff(
		// FIXME splits "strings with whitespaces" into sperate fields
		strings.Fields(spew.Sprintf("%#v", returned)),
		strings.Fields(spew.Sprintf("%#v", want)))

	out := "\n"
	for _, rec := range recs {
		if rec.Delta == difflib.RightOnly {
			out += fmt.Sprintf("%-10s%s\n", "want:", rec.Payload)
		} else if rec.Delta == difflib.LeftOnly {
			out += fmt.Sprintf("%-10s%s\n", "returned:", rec.Payload)
		}
	}
	return out
}

func TestSearchAnnotation(t *testing.T) {

	want := AnnotationSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Annotations: []*Annotation{
			{
				Type:   "release",
				Entity: "bdb24cb5-404b-4f60-bba4-7b730325ae47",
				Name:   "Pieds nus sur la braise",
				Text: `Lyrics and music by Merzhin except:
04, 08, 09, 10 (V. L'hour - Merzhin),
03 (V. L'hour - P. Le Bourdonnec - Merzhin),
05 & 13 (P. Le Bourdonnec - Merzhin),
06 ([http://musicbrainz.org/artist/38cfa519-21bb-4e79-8388-3bf798b8c076.html|JM. Poisson] - Merzhin),
07 ([http://musicbrainz.org/artist/f2d7c07c-a8e7-45c9-a888-0b2e6e3a240d.html|Ignatus] - V. L'hour - Merzhin),
11 ([http://musicbrainz.org/artist/f2d7c07c-a8e7-45c9-a888-0b2e6e3a240d.html|Ignatus] - Merzhin),
12 ([http://musicbrainz.org/artist/38cfa519-21bb-4e79-8388-3bf798b8c076.html|JM. Poisson]).`,
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/annotation", "SearchAnnotation.xml", t)

	returned, err := client.SearchAnnotation("Pieds nus sur la braise", -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Annotations[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}

}

func TestSearchArea(t *testing.T) {

	want := AreaSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Areas: []*Area{
			{
				ID:       "d79e4501-8cba-431b-96e7-bb9976f0ae76",
				Type:     "Subdivision",
				Name:     "Île-de-France",
				SortName: "Île-de-France",
				ISO31662Codes: []ISO31662Code{
					"FR-J",
				},
				Lifespan: Lifespan{
					Ended: false,
				},
				Aliases: []Alias{
					{Locale: "et", SortName: "Île-de-France", Type: "Area name", Primary: "primary", Name: "Île-de-France"},
					{Locale: "ja", SortName: "イル＝ド＝フランス地域圏", Type: "Area name", Primary: "primary", Name: "イル＝ド＝フランス地域圏"},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/area", "SearchArea.xml", t)

	returned, err := client.SearchArea(`"Île-de-France"`, -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Areas[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}

func TestSearchArtist(t *testing.T) {

	want := ArtistSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Artists: []*Artist{
			{
				ID:             "some-artist-id",
				Type:           "Group",
				Name:           "Gopher And Friends",
				Disambiguation: "Some crazy pocket gophers",
				SortName:       "0Gopher And Friends",
				CountryCode:    "DE",
				Lifespan: Lifespan{
					Ended: false,
					Begin: BrainzTime{time.Date(2007, 9, 21, 0, 0, 0, 0, time.UTC)},
					End:   BrainzTime{time.Time{}},
				},
				Aliases: []*Alias{
					{
						Name:     "Mr. Gopher and Friends",
						SortName: "0Mr. Gopher and Friends",
					},
					{
						Name:     "Mr Gopher and Friends",
						SortName: "0Mr Gopher and Friends",
					},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/artist", "SearchArtist.xml", t)

	returned, err := client.SearchArtist("Gopher", -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Artists[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}

func TestSearchRelease(t *testing.T) {

	want := ReleaseSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Releases: []*Release{
			{
				ID:     "9ab1b03e-6722-4ab8-bc7f-a8722f0d34c1",
				Title:  "Fred Schneider & The Shake Society",
				Status: "official",
				TextRepresentation: TextRepresentation{
					Language: "eng",
					Script:   "latn",
				},
				ArtistCredit: ArtistCredit{
					NameCredit{
						Artist{
							ID:       "43bcca8b-9edc-4997-8343-122350e790bf",
							Name:     "Fred Schneider",
							SortName: "Schneider, Fred",
						},
					},
				},
				ReleaseGroup: ReleaseGroup{
					Type: "Album",
				},
				Date:        BrainzTime{time.Date(1991, 4, 30, 0, 0, 0, 0, time.UTC)},
				CountryCode: "us",
				Barcode:     "075992659222",
				Asin:        "075992659222",
				LabelInfos: []*LabelInfo{
					{
						CatalogNumber: "9 26592-2",
						Label: Label{
							Name: "Reprise Records",
						},
					},
				},
				Mediums: []*Medium{
					{
						Format: "cd",
					},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/release", "SearchRelease.xml", t)

	returned, err := client.SearchRelease("Fred", -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Releases[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}

func TestSearchReleaseGroup(t *testing.T) {

	want := ReleaseGroupSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		ReleaseGroups: []*ReleaseGroup{
			{
				ID:          "70664047-2545-4e46-b75f-4556f2a7b83e",
				Type:        "Single",
				Title:       "Main Tenance",
				PrimaryType: "Single",
				ArtistCredit: ArtistCredit{
					NameCredit{
						Artist{
							ID:             "a8fa58d8-f60b-4b83-be7c-aea1af11596b",
							Name:           "Fred Giannelli",
							SortName:       "Giannelli, Fred",
							Disambiguation: "US electronic artist",
						},
					},
				},
				Releases: []*Release{
					{
						ID:    "9168f4cc-a852-4ba5-bf85-602996625651",
						Title: "Main Tenance",
					},
				},
				Tags: []*Tag{
					{
						Count: 1,
						Name:  "electronic",
					},
					{
						Count: 1,
						Name:  "electronica",
					},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/release-group", "SearchReleaseGroup.xml", t)

	returned, err := client.SearchReleaseGroup("Tenance", -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.ReleaseGroups[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}

func TestSearchTag(t *testing.T) {

	want := TagSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  2,
			Offset: 0,
		},
		Tags: []*Tag{
			{
				Name: "shoegaze",
			},
			{
				Name: "rock shoegaze",
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/tag", "SearchTag.xml", t)

	returned, err := client.SearchTag("shoegaze", -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Tags[0]: 100,
		returned.Tags[1]: 62,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}

func TestSearchPlace(t *testing.T) {

	want := PlaceSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Places: []*Place{
			{
				ID:          "d1ab65f8-d082-492a-bd70-ce375548dabf",
				Type:        "Studio",
				Name:        "Chipping Norton Recording Studios",
				Address:     "28–30 New Street, Chipping Norton",
				Coordinates: MBCoordinates{}, // TODO cover
				Area: Area{
					ID:       "44e5e20e-8fbc-4b07-b3f2-22f2199186fd",
					Name:     "Oxfordshire",
					SortName: "Oxfordshire",
				},
				Lifespan: Lifespan{
					Begin: BrainzTime{time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC)},
					End:   BrainzTime{time.Date(1999, 10, 1, 0, 0, 0, 0, time.UTC)},
					Ended: true,
				},
				// TODO Aliases: []*Alias
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/place", "SearchPlace.xml", t)

	returned, err := client.SearchPlace("chipping", -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Places[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
