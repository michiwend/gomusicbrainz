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
	"reflect"
	"testing"
	"time"
)

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
				Status: &Status{
					ID:     "",
					Status: "official",
				},
				TextRepresentation: &TextRepresentation{
					Language: "eng",
					Script:   "latn",
				},
				ArtistCredit: &ArtistCredit{
					NameCredits: []NameCredit{
						{
							Artist:Artist{
								ID:       "43bcca8b-9edc-4997-8343-122350e790bf",
								Name:     "Fred Schneider",
								SortName: "Schneider, Fred",
							},
							JoinPhrase:"",
						},
					},
				},
				ReleaseGroup: &ReleaseGroup{
					Type: "Album",
				},
				Date: &BrainzTime{
					Time:     time.Date(1991, 4, 30, 0, 0, 0, 0, time.UTC),
					Accuracy: Day,
				},
				CountryCode: "us",
				Barcode:     "075992659222",
				Asin:        "075992659222",
				LabelInfos: []LabelInfo{
					{
						CatalogNumber: "9 26592-2",
						Label: &Label{
							Name: "Reprise Records",
						},
					},
				},
				Mediums: []*Medium{
					{
						Format:    Format{Id:"",Name:"cd"},
						TrackList: TrackList{Count: 9},
						DiscList:  DiscList{Count: 2},
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

//TODO implement Lookup test with mediums and tracks
