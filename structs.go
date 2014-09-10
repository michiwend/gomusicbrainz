package gomusicbrainz

type Artist struct {
	Id          string `xml:"id,attr"`
	Type        string `xml:"type,attr"`
	Name        string `xml:"name"`
	SortName    string `xml:"sort-name"`
	CountryCode string `xml:"country"`

	Lifespan struct {
		Ended bool `xml:"ended"`
		//	Begin time.Time `xml:"begin"`
		//	End   time.Time `xml:"end"`
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
