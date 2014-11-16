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
