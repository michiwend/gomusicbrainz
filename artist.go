package gomusicbrainz

// Artist represents generally a musician, a group of musicians, a collaboration
// of multiple musicians or other music professionals.
type Artist struct {
	ID             MBID     `xml:"id,attr"`
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
	Artists []Artist
	Scores  ScoreMap
}

// ResultsWithScore returns a slice of Artists with a specific score.
func (r *ArtistResponse) ResultsWithScore(score int) []Artist {
	var res []Artist
	for k, v := range r.Scores {
		if v == score {
			res = append(res, *k.(*Artist))
		}
	}
	return res
}

type artistListResult struct {
	ArtistList struct {
		WS2ListResponse
		Artists []struct {
			Artist
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"artist"`
	} `xml:"artist-list"`
}
