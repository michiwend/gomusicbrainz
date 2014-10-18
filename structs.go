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

// MBID represents a MusicBrainz ID.
type MBID string

// MBCoordinates represents a tuple of latitude,longitude values.
type MBCoordinates struct {
	// TODO maybe use $geolocation library and its generic type.
	Lat string `xml:"latitude"`
	Lng string `xml:"longitude"`
}

// ScoreMap maps addresses of search request results to its scores.
type ScoreMap map[interface{}]int

type ISO31662Code string

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
	Locale   string `xml:"locale,attr"`
	Type     string `xml:"type,attr"`
	Primary  string `xml:"primary,attr"`
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
