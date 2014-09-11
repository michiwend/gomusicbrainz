/*
 *	Copyright (c) 2014 Michael Wendland
 *
 *	Permission is hereby granted, free of charge, to any person obtaining a
 *	copy of this software and associated documentation files (the "Software"),
 *	to deal in the Software without restriction, including without limitation
 *	the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *	and/or sell copies of the Software, and to permit persons to whom the
 *	Software is furnished to do so, subject to the following conditions:
 *
 *	The above copyright notice and this permission notice shall be included in
 *	all copies or substantial portions of the Software.
 *
 *	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 *	FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 *	IN THE SOFTWARE.
 *
 *	Authors:
 *		Michael Wendland <michael@michiwend.com>
 */

package gomusicbrainz

import (
	"encoding/xml"
	"strings"
	"time"
)

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

type Artist struct {
	Id          string `xml:"id,attr"`
	Type        string `xml:"type,attr"`
	Name        string `xml:"name"`
	SortName    string `xml:"sort-name"`
	CountryCode string `xml:"country"` //ISO_3166-1_alpha-2

	Lifespan struct {
		Ended bool       `xml:"ended"`
		Begin BrainzTime `xml:"begin"`
		End   BrainzTime `xml:"end"`
	} `xml:"life-span"`

	Aliases []struct {
		Name     string `xml:",chardata"`
		SortName string `xml:"sort-name,attr"`
	} `xml:"alias-list"`
}

type ArtistSearchRequest struct {
	ArtistList struct {
		Count   int      `xml:"count,attr"`
		Offset  int      `xml:"offset,attr"`
		Artists []Artist `xml:"artist"`
	} `xml:"artist-list"`
}
