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
				ISO31661Codes: []ISO31661Code{
					"FR",
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
