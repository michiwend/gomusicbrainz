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
				Gender:         "nogender",
				Area: &Area{
					ID:       "some-area-id",
					Name:     "Augsburg",
					SortName: "Augsburg",
				},
				BeginArea: &Area{
					ID:       "some-area-id",
					Name:     "Mountain View",
					SortName: "Mountain View",
				},
				Lifespan: &Lifespan{
					Ended: false,
					Begin: &BrainzTime{
						Time:     time.Date(2007, 9, 21, 0, 0, 0, 0, time.UTC),
						Accuracy: Day,
					},
					End: &BrainzTime{Time: time.Time{}},
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
				Tags: []Tag{
					{
						Count: 1,
						Name:  "Pocket Gopher Music",
					},
					{
						Count: 2,
						Name:  "Golang",
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

func TestLookupArtist(t *testing.T) {

	want := Artist{
		ID:             "10adbe5e-a2c0-4bf3-8249-2b4cbf6e6ca8",
		Type:           "Group",
		Name:           "Massive Attack",
		Disambiguation: "",
		SortName:       "Massive Attack",
		CountryCode:    "",
		Area: &Area{
			ID:       "40d758a4-b7c2-40f3-b439-5efbd2a3b038",
			Name:     "Bristol",
			SortName: "Bristol",
			ISO31661Codes: []ISO31661Code{
				"GB",
			},
		},
		BeginArea: &Area{
			ID:       "40d758a4-b7c2-40f3-b439-5efbd2a3b038",
			Name:     "Bristol",
			SortName: "Bristol",
			ISO31661Codes: []ISO31661Code{
				"GB",
			},
		},
		Lifespan: &Lifespan{
			Ended: false,
			Begin: &BrainzTime{
				Time:     time.Date(1987, 1, 1, 0, 0, 0, 0, time.UTC),
				Accuracy: Year,
			},
			End: &BrainzTime{Time: time.Time{}},
		},
		Relations: TargetRelationsMap{
			"artist": []Relation{
				&ArtistRelation{
					RelationAbstract: RelationAbstract{
						TypeID:    "5be4c609-9afa-4ea0-910b-12ffb71e3821",
						Type:      "member of band",
						Target:    "54912e02-166c-49fe-ba95-cd77ef182390",
						Direction: "backward",
						Begin: BrainzTime{
							Time:     time.Date(1987, 1, 1, 0, 0, 0, 0, time.UTC),
							Accuracy: Year,
						},
						End: BrainzTime{
							Time:     time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC),
							Accuracy: Year,
						},
						Ended: true,
					},
					// TODO Attribute list
					Artist: Artist{
						ID:             "54912e02-166c-49fe-ba95-cd77ef182390",
						Name:           "Mushroom",
						SortName:       "Mushroom",
						Disambiguation: "Andrew Vowles, member of Massive Attack",
					},
				},
			},
			"release": []Relation{
				&ReleaseRelation{
					RelationAbstract: RelationAbstract{
						TypeID: "307e95dd-88b5-419b-8223-b146d4a0d439",
						Type:   "design/illustration",
						Target: "07832b54-8266-47d5-bb0e-62c7f2cf5da5",
					},
					Release: Release{
						ID:      "07832b54-8266-47d5-bb0e-62c7f2cf5da5",
						Title:   "Protection",
						Quality: "normal",
						Date: &BrainzTime{
							Time:     time.Date(1995, 1, 24, 0, 0, 0, 0, time.UTC),
							Accuracy: Day,
						},
						CountryCode: "US",
						Barcode:     "724383988327",
					},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile(
		"/artist/10adbe5e-a2c0-4bf3-8249-2b4cbf6e6ca8",
		"LookupArtist.xml", t)

	returned, err := client.LookupArtist(
		"10adbe5e-a2c0-4bf3-8249-2b4cbf6e6ca8",
		"artist-rels",
		"release-rels")

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}

}
