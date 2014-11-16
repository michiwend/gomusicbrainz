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
	"path"
	"testing"

	"github.com/michiwend/golang-pretty"
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

// serveTestFile responses to the http client with content of a test file
// located in ./testdata
func serveTestFile(endpoint string, testfile string, t *testing.T) {

	//TODO check request URL if it matches one of the following patterns
	//lookup:   /<ENTITY>/<MBID>?inc=<INC>
	//browse:   /<ENTITY>?<ENTITY>=<MBID>&limit=<LIMIT>&offset=<OFFSET>&inc=<INC>
	//search:   /<ENTITY>?query=<QUERY>&limit=<LIMIT>&offset=<OFFSET>

	t.Log("Handling endpoint", endpoint)
	t.Log("Serving test file", testfile)

	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {

		t.Log("GET request was:", r.URL.String())

		http.ServeFile(w, r, path.Join("./testdata", testfile))
	})
}

// pretty prints a diff
func requestDiff(want, returned interface{}) string {

	out := "\n"

	for _, diff := range pretty.Diff(want, returned) {
		out += fmt.Sprintln("difference in", diff)
	}
	return out
}
