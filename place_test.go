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
					Begin: &BrainzTime{
						Time:     time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC),
						Accuracy: Year,
					},
					End: &BrainzTime{
						Time:     time.Date(1999, 10, 1, 0, 0, 0, 0, time.UTC),
						Accuracy: Month,
					},
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
