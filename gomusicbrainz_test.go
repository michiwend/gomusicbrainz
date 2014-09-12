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

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client GoMusicBrainz
)

// Init multiplexer and httptest server
func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	host, _ := url.Parse(server.URL)
	client = GoMusicBrainz{WS2RootURL: host}
}

// The handleFunc simply passes response to the http client.
func handleFunc(url string, response *string, t *testing.T) {
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, *response)
	})
}

func TestSearchArtist(t *testing.T) {

	setup()
	defer server.Close()

	response := `
		<?xml version="1.0" standalone="yes"?>
		<metadata xmlns="http://musicbrainz.org/ns/mmd-2.0#" xmlns:ext="http://musicbrainz.org/ns/ext#-2.0" created="2014-09-12T06:31:24.904Z">
			<artist-list count="1" offset="0">
				<artist id="some-artist-id" type="Group" ext:score="100">
					<name>Gopher And Friends</name>
					<sort-name>0Gopher And Friends</sort-name>
					<country>DE</country>
					<area id="some-area-id">
						<name>Augsburg</name>
						<sort-name>Augsburg</sort-name>
					</area>
					<begin-area id="some-area-id">
						<name>Mountain View</name>
						<sort-name>Mountain View</sort-name>
					</begin-area>
					<disambiguation>Some crazy pocket gophers</disambiguation>
					<life-span>
						<begin>2007-09-21</begin>
						<ended>false</ended>
					</life-span>
					<alias-list>
						<alias sort-name="0Mr. Gopher and Friends">Mr. Gopher and Friends</alias>
						<alias sort-name="0Mr Gopher and Friends">Mr Gopher and Friends</alias>
					</alias-list>
					<tag-list>
						<tag count="1">
							<name>Pocket Gopher Music</name>
						</tag>
						<tag count="2">
							<name>Golang</name>
						</tag>
					</tag-list>
				</artist>
			</artist-list>
		</metadata>`

	want := []Artist{
		{
			Id:          "some-artist-id",
			Type:        "Group",
			Name:        "Gopher And Friends",
			SortName:    "0Gopher And Friends",
			CountryCode: "DE",
			Lifespan: Lifespan{
				Ended: false,
				Begin: BrainzTime{time.Date(2007, 9, 21, 0, 0, 0, 0, time.UTC)},
				End:   BrainzTime{time.Time{}},
			},
			Aliases: []Alias{
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
	}

	handleFunc("/artist", &response, t)

	returned, err := client.SearchArtist("", -1, -1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(returned, want) {
		t.Errorf("Artists returned: %+v, want: %+v", returned, want)
	}
}

func TestSearchRelease(t *testing.T) {

	setup()
	defer server.Close()

	response := `
		<metadata xmlns="http://musicbrainz.org/ns/mmd-2.0#" xmlns:ext="http://musicbrainz.org/ns/ext#-2.0">
			<release-list offset="0" count="1">
				<release id="9ab1b03e-6722-4ab8-bc7f-a8722f0d34c1" ext:score="100">
					<title>Fred Schneider &amp; The Shake Society</title>
					<status>official</status>
					<text-representation>
						<language>eng</language>
						<script>latn</script>
					</text-representation>
					<artist-credit>
						<name-credit>
							<artist id="43bcca8b-9edc-4997-8343-122350e790bf">
							   <name>Fred Schneider</name>
							   <sort-name>Schneider, Fred</sort-name>
							</artist>
						</name-credit>
					</artist-credit>
					<release-group type="Album"/>
					<date>1991-04-30</date>
					<country>us</country>
					<barcode>075992659222</barcode>
					<asin>075992659222</asin>
					<label-info-list>
						<label-info>
							<catalog-number>9 26592-2</catalog-number>
							<label>
								<name>Reprise Records</name>
							</label>
						</label-info>
					</label-info-list>
					<medium-list>
						<medium><format>cd</format>
							<disc-list count="2"/>
							<track-list count="9"/>
						 </medium>
					</medium-list>
				</release>
			</release-list>
		</metadata>`

	want := []Release{
		{
			Id:     "9ab1b03e-6722-4ab8-bc7f-a8722f0d34c1",
			Title:  "Fred Schneider & The Shake Society",
			Status: "official",
			TextRepresentation: TextRepresentation{
				Language: "eng",
				Script:   "latn",
			},
			ArtistCredit: ArtistCredit{
				NameCredit{
					Artist{
						Id:       "43bcca8b-9edc-4997-8343-122350e790bf",
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
			LabelInfos: []LabelInfo{
				{
					CatalogNumber: "9 26592-2",
					Label: Label{
						Name: "Reprise Records",
					},
				},
			},
			Mediums: []Medium{
				{
					Format: "cd",
				},
			},
		},
	}

	handleFunc("/release", &response, t)

	returned, err := client.SearchRelease("", -1, -1)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(returned, want) {
		t.Errorf("Releases returned: %+v, want: %+v", returned, want)
	}
}
