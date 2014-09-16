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

package gomusicbrainz

import (
	"encoding/xml"
	"strings"
	"time"
)

// BrainzTime implements XMLUnmarshaler interface and is used to unmarshal the
// XML date fields.
type BrainzTime struct {
	time.Time
}

func (t *BrainzTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	var p time.Time
	var err error
	d.DecodeElement(&v, &start)

	switch strings.Count(v, "-") {
	case 0:
		p, err = time.Parse("2006", v)
	case 1:
		p, err = time.Parse("2006-01", v)
	case 2:
		p, err = time.Parse("2006-01-02", v)
	}

	// TODO handle empty fields

	if err != nil {
		return err
	}
	*t = BrainzTime{p}
	return nil
}

// WS2ListResponse is an abstract common type that provides the Count and Offset
// fields for ervery list response.
type WS2ListResponse struct {
	Count  int `xml:"count,attr"`
	Offset int `xml:"offset,attr"`
}

type Lifespan struct {
	Ended bool       `xml:"ended"`
	Begin BrainzTime `xml:"begin"`
	End   BrainzTime `xml:"end"`
}

// Alias is a type for aliases/misspellings of artists, works, areas, labels
// and places.
type Alias struct {
	Name     string `xml:",chardata"`
	SortName string `xml:"sort-name,attr"`
}

// Tag is the common type for Tags.
type Tag struct {
	Count int    `xml:"count,attr"`
	Name  string `xml:"name"`
	Score int    `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
}

// TagResponse is the response type returned by tag request methods.
type TagResponse struct {
	WS2ListResponse
	Tags []Tag `xml:"tag"`
}

type Artist struct {
	Id             string   `xml:"id,attr"`
	Type           string   `xml:"type,attr"`
	Name           string   `xml:"name"`
	Disambiguation string   `xml:"disambiguation"`
	SortName       string   `xml:"sort-name"`
	CountryCode    string   `xml:"country"`
	Lifespan       Lifespan `xml:"life-span"`
	Aliases        []Alias  `xml:"alias-list>alias"`
}

// ArtistResponse is the response type returned by artist request methods.
type ArtistResponse struct {
	WS2ListResponse
	Artists []Artist `xml:"artist"`
}

type Label struct {
	Name string `xml:"name"`
}

type LabelResponse struct {
	//TODO implement
}

type LabelInfo struct {
	CatalogNumber string `xml:"catalog-number"`
	Label         Label  `xml:"label"`
}

type Medium struct {
	Format string `xml:"format"`
	//DiscList TODO implement type
	//TrackList TODO implement type
}

type TextRepresentation struct {
	Language string `xml:"language"`
	Script   string `xml:"script"`
}

type ArtistCredit struct {
	NameCredit NameCredit `xml:"name-credit"`
}

type NameCredit struct {
	Artist Artist `xml:"artist"`
}

type Release struct {
	Id                 string             `xml:"id,attr"`
	Title              string             `xml:"title"`
	Status             string             `xml:"status"`
	Disambiguation     string             `xml:"disambiguation"`
	TextRepresentation TextRepresentation `xml:"text-representation"`
	ArtistCredit       ArtistCredit       `xml:"artist-credit"`
	ReleaseGroup       ReleaseGroup       `xml:"release-group"`
	Date               BrainzTime         `xml:"date"`
	CountryCode        string             `xml:"country"`
	Barcode            string             `xml:"barcode"`
	Asin               string             `xml:"asin"`
	LabelInfos         []LabelInfo        `xml:"label-info-list>label-info"`
	Mediums            []Medium           `xml:"medium-list>medium"`
}

// ReleaseResponse is the response type returned by release request methods.
type ReleaseResponse struct {
	WS2ListResponse
	Releases []Release `xml:"release"`
}

type ReleaseGroup struct {
	Id           string       `xml:"id,attr"`
	Type         string       `xml:"type,attr"`
	PrimaryType  string       `xml:"primary-type"`
	Title        string       `xml:"title"`
	ArtistCredit ArtistCredit `xml:"artist-credit"`
	Releases     []Release    `xml:"release-list>release"` // FIXME if important unmarshal count,attr
	Tags         []Tag        `xml:"tag-list>tag"`
}

// ReleaseGroupResponse is the response type returned by release group
// request methods.
type ReleaseGroupResponse struct {
	WS2ListResponse
	ReleaseGroups []ReleaseGroup `xml:"release-group"`
}

// Annotation is a miniature wiki that can be added to any existing artists,
// labels, recordings, releases, release groups and works.
type Annotation struct {
	Type   string `xml:"type,attr"`
	Entity string `xml:"entity"`
	Name   string `xml:"name"`
	Text   string `xml:"text"`
}

// AnnotaionResponse is the response type returned by annotation request
// methods.
type AnnotationResponse struct {
	WS2ListResponse
	Annotations []Annotation `xml:"annotation"`
}

type Area struct {
	//TODO implement
}

type AreaResponse struct {
	//TODO implement
}
type CDStubs struct {
	//TODO implement
}

type CDStubsResponse struct {
	//TODO implement
}
type Freedb struct {
	//TODO implement
}

type FreedbResponse struct {
	//TODO implement
}

type Place struct {
	//TODO implement
}

type PlaceResponse struct {
	//TODO implement
}
type Recording struct {
	//TODO implement
}

type RecordingResponse struct {
	//TODO implement
}
type Work struct {
	//TODO implement
}

type WorkResponse struct {
	//TODO implement
}
