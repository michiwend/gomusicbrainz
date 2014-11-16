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

func TestSearchCDStub(t *testing.T) {

	want := CDStubSearchResponse{
		WS2ListResponse: WS2ListResponse{
			Count:  1,
			Offset: 0,
		},
		CDStubs: []*CDStub{
			{
				ID:      "vi44WFVS5zRT2svM.PORcEm9LJk-",
				Title:   "Silent Conflict (Live @ The Hard Rock Cafe)",
				Artist:  "Bonobo",
				Barcode: "634479355059",
				Comment: "CD Baby id:bonobo",
				TrackList: struct {
					Count int `xml:"count,attr"`
				}{
					Count: 3,
				},
			},
		},
	}

	setupHTTPTesting()
	defer server.Close()
	serveTestFile("/cdstub", "SearchCDStub.xml", t)

	returned, err := client.SearchCDStub(`bonobo`, -1, -1)
	if err != nil {
		t.Error(err)
	}

	want.Scores = ScoreMap{
		returned.CDStubs[0]: 100,
	}

	if !reflect.DeepEqual(*returned, want) {
		t.Error(requestDiff(&want, returned))
	}
}
