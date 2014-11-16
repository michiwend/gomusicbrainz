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

func TestSearchLabel(t *testing.T) {

	want := LabelSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		Labels: []*Label{
			{
				ID:             "c1c625b5-9929-4a30-8c3e-f77e109cdf07",
				Type:           "Original Production",
				Name:           "Compost Records",
				SortName:       "Compost Records",
				Disambiguation: "German record label established in 1994.",
				CountryCode:    "DE",
				LabelCode:      2518,
				Area: Area{
					ID:       "85752fda-13c4-31a3-bee5-0e5cb1f51dad",
					Name:     "Germany",
					SortName: "Germany",
				},
				Lifespan: Lifespan{
					Begin: BrainzTime{time.Date(1994, 1, 1, 0, 0, 0, 0, time.UTC)},
					Ended: false,
				},
				Aliases: []*Alias{
					{
						Locale:   "ja",
						SortName: "コンポスト・レコーズ",
						Name:     "コンポスト・レコーズ",
						Type:     "Label name",
					},
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/label", "SearchLabel.xml", t)

	returned, err := client.SearchLabel(`label:"Compost%20Records"`, -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.Labels[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
