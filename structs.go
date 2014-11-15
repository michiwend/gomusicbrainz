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

// MBLookupEntity represents all MusicBrainz entities for which a MBID-lookup
// request is provided by WS2.
type MBLookupEntity interface {
	lookupResult() interface{}
	apiEndpoint() string
	id() MBID
}

type MBEntity interface {
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

// Lifespan represents either the life span of a natural person or more
// generally the period of time in which an entity e.g. a Label existed.
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

// Medium represents one of the physical, separate things you would get when
// you buy something in a record store e.g. CDs, vinyls, etc. Mediums are
// always included in a release. For more information visit
// https://musicbrainz.org/doc/Medium
type Medium struct {
	Format string `xml:"format"`
	//DiscList TODO implement type
	//TrackList TODO implement type
}

type TextRepresentation struct {
	Language string `xml:"language"`
	Script   string `xml:"script"`
}

// ArtistCredit is either used to link multiple artists to one
// release/recording or to credit an artist with a different name.
// Visist https://musicbrainz.org/doc/Artist_Credit for more information.
type ArtistCredit struct {
	NameCredit NameCredit `xml:"name-credit"`
}

type NameCredit struct {
	Artist Artist `xml:"artist"`
}

type Relation interface {
}

type URLRelation struct {
	RelationAbstract
}

// RelationAbstract is the common abstract type for Relations.
type RelationAbstract struct {
	TypeID MBID   `xml:"type-id,attr"`
	Type   string `xml:"type,attr"`
	Target MBID   `xml:"target"`
}

// ReleaseRelation is the Relation type for Releases.
type ReleaseRelation struct {
	RelationAbstract
	Release Release `xml:"release"`
}

// ArtistRelation is the Relation type for Artists.
type ArtistRelation struct {
	RelationAbstract
	Artist    Artist `xml:"artist"`
	Direction string `xml:"direction"`
}

// TargetRelationsMap maps target-types to Relations.
type TargetRelationsMap map[string][]Relation

// UnmarshalXML is needed to implement XMLUnmarshaler for custom, value-based
// unmarshaling of relation-list elements.
func (r *TargetRelationsMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var targetType string

	for _, v := range start.Attr {
		if v.Name.Local == "target-type" {
			targetType = v.Value
			break
		}
	}

	if *r == nil {
		(*r) = make(map[string][]Relation)
	}

	switch targetType {
	case "artist":
		var res struct {
			XMLName   xml.Name          `xml:"relation-list"`
			Relations []*ArtistRelation `xml:"relation"`
		}
		if err := d.DecodeElement(&res, &start); err != nil {
			return err
		}

		(*r)[targetType] = make([]Relation, len(res.Relations))

		for i, v := range res.Relations {
			(*r)[targetType][i] = v
		}

	case "release":
		var res struct {
			XMLName   xml.Name           `xml:"relation-list"`
			Relations []*ReleaseRelation `xml:"relation"`
		}

		if err := d.DecodeElement(&res, &start); err != nil {
			return err
		}

		(*r)[targetType] = make([]Relation, len(res.Relations))

		for i, v := range res.Relations {
			(*r)[targetType][i] = v
		}

	// FIXME implement missing relations

	default:
		return d.Skip()
	}

	return nil
}
