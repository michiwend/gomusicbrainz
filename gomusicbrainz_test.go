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

// TODO use testdata from https://github.com/metabrainz/mmd-schema/tree/master/test-data/valid

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aryann/difflib"
	"github.com/davecgh/go-spew/spew"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *WS2Client
)

// Init multiplexer and httptest server
func setupHTTPTesting() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client, _ = NewWS2Client(
		server.URL,
		"Application Name",
		"Version",
		"Contact",
	)

	// NOTE this fixes testing since the test server does not listen on /ws/2
	client.WS2RootURL.Path = ""
}

// handleFunc passes response to the http client.
/*
func handleFunc(url string, response *string, t *testing.T) {
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, *response)
	})
}*/

// serveTestFile responses to the http client with content of a test file
// located in ./testdata
func serveTestFile(url string, testfile string, t *testing.T) {

	//TODO check request URL if it matches one of the following patterns
	//lookup:   /<ENTITY>/<MBID>?inc=<INC>
	//browse:   /<ENTITY>?<ENTITY>=<MBID>&limit=<LIMIT>&offset=<OFFSET>&inc=<INC>
	//search:   /<ENTITY>?query=<QUERY>&limit=<LIMIT>&offset=<OFFSET>

	t.Log("Listening on", url)

	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		if testing.Verbose() {
			fmt.Println("GET request:", r.URL.String())
		}

		http.ServeFile(w, r, "./testdata/"+testfile)
	})
}

// pretty prints a diff
func requestDiff(want, returned interface{}) string {
	spew.Config.SortKeys = true
	spew.Config.ContinueOnMethod = true
	recs := difflib.Diff(
		// FIXME splits "strings with whitespaces" into sperate fields
		strings.Fields(spew.Sprintf("%#v", returned)),
		strings.Fields(spew.Sprintf("%#v", want)))

	out := "\n"
	for _, rec := range recs {
		if rec.Delta == difflib.RightOnly {
			out += fmt.Sprintf("%-10s%s\n", "want:", rec.Payload)
		} else if rec.Delta == difflib.LeftOnly {
			out += fmt.Sprintf("%-10s%s\n", "returned:", rec.Payload)
		}
	}
	return out

	/*
		out := "\n"

		for _, diff := range pretty.Diff(want, returned) {
			out += fmt.Sprintln("difference in", diff)
		}
		return out
	*/

}
