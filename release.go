package gomusicbrainz

// Release represents a unique release (i.e. issuing) of a product on a
// specific date with specific release information such as the country, label,
// barcode, packaging, etc. More information at https://musicbrainz.org/doc/Release
type Release struct {
	ID                 string             `xml:"id,attr"`
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
	Releases []Release
	Scores   ScoreMap
}

type releaseListResult struct {
	ReleaseList struct {
		WS2ListResponse
		Releases []struct {
			Release
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"release"`
	} `xml:"release-list"`
}
