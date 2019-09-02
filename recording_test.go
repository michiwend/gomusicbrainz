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

func TestSearchRecording(t *testing.T) {

	want := RecordingSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Recordings: []*Recording{
			{
				ID:     "07339604-c19c-4efe-9195-f9c3b127a458",
				Title:  "Fred",
				Length: 473000,
				ArtistCredit: ArtistCredit{
					NameCredits: []NameCredit{
						NameCredit{
							Artist{
								ID:       "695e75b5-c6db-43ee-abeb-2f3e50d96c3e",
								Name:     "Imperiet",
								SortName: "Imperiet",
							},
						},
					},
				},
				IsrcList: []Isrc{{Id: "USDJ20400074"}},

				//TODO add missing fields
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/recording", "SearchRecording.xml", t)

	returned, err := client.SearchRecording("Fred", -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Recordings[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
