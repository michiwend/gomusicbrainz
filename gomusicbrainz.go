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
 * 	Authors:
 * 		Michael Wendland <michael@michiwend.com>
 */

/*
Package gomusicbrainz implements a MusicBrainz WS2 client library.

MusicBrainz WS2 (Version 2 of the XML Web Service) supports three different requests:

Search requests

With search requests you can search MusicBrainz´ database for all entities.
GoMusicBrainz implements one search method for every search request in the form:

	func (*WS2Client) Search<TYPE>(searchTerm, limit, offset) (Response<TYPE>, error)

searchTerm follows the Apache Lucene syntax and can either contain multiple
fields with logical operators or just a simple search string. Please refer to
https://lucene.apache.org/core/4_3_0/queryparser/org/apache/lucene/queryparser/classic/package-summary.html#package_description
for more details on the lucene syntax. limit defines how many entries should be
returned (1-100, default 25). offset is used for paging through more than one
page of results. To ignore limit and/or offset, set it to -1.

Lookup requests

TODO

Browse requets

TODO

*/
package gomusicbrainz

import (
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// NewWS2Client returns a new instance of WS2Client. Please provide meaningful
// information about your application as described at
// https://musicbrainz.org/doc/XML_Web_Service/Rate_Limiting#Provide_meaningful_User-Agent_strings
func NewWS2Client(rooturl string, appname string, version string, contact string) *WS2Client {
	c := WS2Client{}

	c.WS2RootURL, _ = url.Parse(rooturl)
	c.userAgentHeader = appname + "/" + version + " ( " + contact + " ) "

	return &c
}

// WS2Client defines a Go client for the MusicBrainz Web Service 2.
type WS2Client struct {
	WS2RootURL *url.URL // The API root URL

	userAgentHeader string
}

func (c *WS2Client) getReqeust(data interface{}, params url.Values, endpoint string) error {

	client := &http.Client{}

	req, err := http.NewRequest("GET", c.WS2RootURL.String()+endpoint+"?"+params.Encode(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", c.userAgentHeader)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)

	if err = decoder.Decode(data); err != nil {
		return err
	}
	return nil

}

// intParamToString returns an empty string for -1.
func intParamToString(i int) string {
	if i == -1 {
		return ""
	}
	return strconv.Itoa(i)
}

func (c *WS2Client) searchRequest(endpoint string, result interface{}, searchTerm string, limit, offset int) error {

	params := url.Values{
		"query":  {searchTerm},
		"limit":  {intParamToString(limit)},
		"offset": {intParamToString(offset)},
	}

	if err := c.getReqeust(result, params, endpoint); err != nil {
		return err
	}

	return nil
}

// SetRootURL sets the root URL for WS2.
func (c *WS2Client) SetRootURL(rooturl string) error {
	var err error
	c.WS2RootURL, err = url.Parse(rooturl)
	return err
}

// SetClientInfo sets the HTTP user-agent header of the WS2Client. Please
// provide meaningful information about your application as described at
// https://musicbrainz.org/doc/XML_Web_Service/Rate_Limiting#Provide_meaningful_User-Agent_strings
func (c *WS2Client) SetClientInfo(application string, version string, contact string) {
	c.userAgentHeader = application + "/" + version + " ( " + contact + " ) "
}

// SearchAnnotation queries MusicBrainz´ Search Server for Annotations.
// With no fields specified searchTerm searches TODO. For a list of all valid
// search fields visit
// http://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Annotation
func (c *WS2Client) SearchAnnotation(searchTerm string, limit, offset int) (*AnnotationResponse, error) {

	var result struct {
		Response AnnotationResponse `xml:"annotation-list"`
	}

	err := c.searchRequest("/annotation", &result, searchTerm, limit, offset)

	return &result.Response, err
}

// SearchArea queries MusicBrainz´ Search Server for Areas.
// With no fields specified searchTerm searches the area and sortname fields.
// For a list of all valid search fields visit
// http://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Area
func (c *WS2Client) SearchArea(searchTerm string, limit, offset int) (*AreaResponse, error) {

	var result struct {
		Response AreaResponse `xml:"area-list"`
	}

	err := c.searchRequest("/area", &result, searchTerm, limit, offset)

	return &result.Response, err
}

// SearchArtist queries MusicBrainz´ Search Server for Artists.
// With no fields specified searchTerm searches the artist, sortname and alias
// fields. For a list of all valid fields visit
// http://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Artist
func (c *WS2Client) SearchArtist(searchTerm string, limit, offset int) (*ArtistResponse, error) {

	var result struct {
		Response ArtistResponse `xml:"artist-list"`
	}

	err := c.searchRequest("/artist", &result, searchTerm, limit, offset)

	return &result.Response, err
}

// SearchRelease queries MusicBrainz´ Search Server for Releases.
// With no fields specified searchTerm searches the release field only. For a
// list of all valid fields visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Release
func (c *WS2Client) SearchRelease(searchTerm string, limit, offset int) (*ReleaseResponse, error) {

	var result struct {
		Response ReleaseResponse `xml:"release-list"`
	}

	err := c.searchRequest("/release", &result, searchTerm, limit, offset)

	return &result.Response, err
}

// SearchReleaseGroup queries MusicBrainz´ Search Server for ReleaseGroups.
// With no fields specified searchTerm searches the releasgroup field only. For
// a list of all valid fields visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Release_Group
func (c *WS2Client) SearchReleaseGroup(searchTerm string, limit, offset int) (*ReleaseGroupResponse, error) {

	var result struct {
		Response ReleaseGroupResponse `xml:"release-group-list"`
	}

	err := c.searchRequest("/release-group", &result, searchTerm, limit, offset)

	return &result.Response, err
}

// SearchTag queries MusicBrainz' Search Server for Tags.
// searchTerm only contains the tag field. For more information visit
// https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search#Tag
func (c *WS2Client) SearchTag(searchTerm string, limit, offset int) (*TagResponse, error) {

	var result struct {
		Response TagResponse `xml:"tag-list"`
	}

	err := c.searchRequest("/tag", &result, searchTerm, limit, offset)

	return &result.Response, err
}

func (c *WS2Client) SearchCDStubs(searchTerm string, limit, offset int) (*CDStubsResponse, error) {
	//TODO implement
	return nil, nil
}
func (c *WS2Client) SearchFreedb(searchTerm string, limit, offset int) (*FreedbResponse, error) {
	//TODO implement
	return nil, nil
}
func (c *WS2Client) SearchLabel(searchTerm string, limit, offset int) (*LabelResponse, error) {
	//TODO implement
	return nil, nil
}
func (c *WS2Client) SearchPlace(searchTerm string, limit, offset int) (*PlaceResponse, error) {
	//TODO implement
	return nil, nil
}
func (c *WS2Client) SearchRecording(searchTerm string, limit, offset int) (*RecordingResponse, error) {
	//TODO implement
	return nil, nil
}
func (c *WS2Client) SearchWork(searchTerm string, limit, offset int) (*WorkResponse, error) {
	//TODO implement
	return nil, nil
}
