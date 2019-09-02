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

// MBID represents a MusicBrainz Identifier. A MBID is a 36 character
// Universally Unique Identifier that is permanently assigned to each entity in
// the database, i.e. artists, release groups, releases, recordings, works,
// labels, areas, places and URLs.
type MBID string

// MBentity is an interface implemented by all MusicBrainz entities with MBIDs.
type MBEntity interface {
	Id() MBID
	apiEndpoint() string
}

// MBLookupEntity represents all MusicBrainz entities for which a MBID-lookup
// request is provided by WS2.
type MBLookupEntity interface {
	MBEntity
	lookupResult() interface{}
}

// MBCoordinates represents a tuple of latitude,longitude values.
type MBCoordinates struct {
	Lat string `xml:"latitude"`
	Lng string `xml:"longitude"`
}

// ScoreMap maps addresses of search request results to its scores.
type ScoreMap map[interface{}]int

type ISO31662Code string

// BrainzTimeAccuracy specifies the accuracy for the corresponding BrainzTime.
type BrainzTimeAccuracy int

const (
	Year BrainzTimeAccuracy = iota
	Month
	Day
)

// BrainzTime represents a MusicBrainz date by combining time.Time with a
// Accuracy field to distinguish between different date accuracies e.g. "2006"
// and "2006-01".
//
// You can compare the accuracy of a BrainzTime type simply by using operators:
//
//	// time1 represents "2006-01"
//	// time2 represents "2006"
//
//	time1.Accuracy > time2.Accuray // true
type BrainzTime struct {
	time.Time
	Accuracy BrainzTimeAccuracy
}

func (t *BrainzTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	var err error
	d.DecodeElement(&v, &start)

	if v != "" {
		switch strings.Count(v, "-") {
		case 0:
			t.Time, err = time.Parse("2006", v)
			t.Accuracy = Year
		case 1:
			t.Time, err = time.Parse("2006-01", v)
			t.Accuracy = Month
		case 2:
			t.Time, err = time.Parse("2006-01-02", v)
			t.Accuracy = Day
		}
	}

	return err
}

// WS2ListResponse is a abstract common type that provides the Count and Offset
// fields for ervery list response.
type WS2ListResponse struct {
	Count  int `xml:"count,attr"`
	Offset int `xml:"offset,attr"`
}

// Lifespan represents either the life span of a natural person or more
// generally the period of time in which an entity e.g. a Label existed.
type Lifespan struct {
	Begin BrainzTime `xml:"begin"`
	End   BrainzTime `xml:"end"`
	Ended bool       `xml:"ended"`
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

type TrackList struct {
	Count  int      `xml:"count,attr"`
	Tracks []*Track `xml:"track"`
}

type DiscList struct {
	Count int `xml:"count,attr"`
}

// Medium represents one of the physical, separate things you would get when
// you buy something in a record store e.g. CDs, vinyls, etc. Mediums are
// always included in a release. For more information visit
// https://musicbrainz.org/doc/Medium
type Medium struct {
	Format    string    `xml:"format"`
	Position  int       `xml:"position"`
	DiscList  DiscList  `xml:"disc-list"`
	TrackList TrackList `xml:"track-list"`
}

// Track represents a recording on a particular release (or, more exactly, on
// a particular medium). See https://musicbrainz.org/doc/Track
type Track struct {
	ID        MBID      `xml:"id,attr"`
	Position  int       `xml:"position"`
	Number    string    `xml:"number"`
	Length    int       `xml:"length"`
	Recording Recording `xml:"recording"`
}

type TextRepresentation struct {
	Language string `xml:"language"`
	Script   string `xml:"script"`
}

// ArtistCredit is either used to link multiple artists to one
// release/recording or to credit an artist with a different name.
// Visist https://musicbrainz.org/doc/Artist_Credit for more information.
type ArtistCredit struct {
	NameCredits []NameCredit `xml:"name-credit"`
}

type NameCredit struct {
	Artist Artist `xml:"artist"`
}

type Isrc struct {
	Id string `xml:"id,attr"`
}

// Relation describes a relationship between different MusicBrainz entities.
// See this link https://musicbrainz.org/relationships for a table of
// relationships.
type Relation interface {
	TypeOf() string
}

// RelationAbstract is the common abstract type for Relations.
type RelationAbstract struct {
	Type        string     `xml:"type,attr"`
	TypeID      MBID       `xml:"type-id,attr"`
	Target      string     `xml:"target"`
	TargetID    MBID       `xml:"target-id,attr"`
	OrderingKey int        `xml:"ordering-key"`
	Direction   string     `xml:"direction"`
	Begin       BrainzTime `xml:"begin"`
	End         BrainzTime `xml:"end"`
	Ended       bool       `xml:"ended"`
}

func (r *RelationAbstract) TypeOf() string {
	return r.Type
}

// RelationsOfTypes returns a slice of Relations for the given relTypes. For a
// list of all possible relationships see https://musicbrainz.org/relationships
func RelationsOfTypes(rels []Relation, relTypes ...string) []Relation {

	var out []Relation

	for _, rel := range rels {
		for _, relType := range relTypes {
			if rel.TypeOf() == relType {
				out = append(out, rel)
			}
		}
	}

	return out
}

type URLRelation struct {
	RelationAbstract
}

// ReleaseRelation is the Relation type for Releases.
type ReleaseRelation struct {
	RelationAbstract
	Release Release `xml:"release"`
}

// ArtistRelation is the Relation type for Artists.
type ArtistRelation struct {
	RelationAbstract
	Artist Artist `xml:"artist"`
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

	case "url":
		var res struct {
			XMLName   xml.Name       `xml:"relation-list"`
			Relations []*URLRelation `xml:"relation"`
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
