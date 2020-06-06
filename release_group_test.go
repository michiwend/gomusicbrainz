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
)

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
					NameCredits: []NameCredit{
						NameCredit{
							Artist{
								ID:             "a8fa58d8-f60b-4b83-be7c-aea1af11596b",
								Name:           "Fred Giannelli",
								SortName:       "Giannelli, Fred",
								Disambiguation: "US electronic artist",
							},
							"",
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
