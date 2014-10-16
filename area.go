package gomusicbrainz

// Area represents a geographic region or settlement.
type Area struct {
	ID            string         `xml:"id,attr"`
	Type          string         `xml:"type,attr"`
	Name          string         `xml:"name"`
	SortName      string         `xml:"sort-name"`
	ISO31662Codes []ISO31662Code `xml:"iso-3166-2-code-list>iso-3166-2-code"`
	Lifespan      Lifespan       `xml:"life-span"`
	Aliases       []Alias        `xml:"alias-list>alias"`
}

// AreaResponse is the response type returned by area request methods.
type AreaResponse struct {
	WS2ListResponse
	Areas  []Area `xml:"area"`
	Scores ScoreMap
}

// ResultsWithScore returns a slice of Areas with a specific score.
func (r *AreaResponse) ResultsWithScore(score int) []Area {
	var res []Area
	for k, v := range r.Scores {
		if v == score {
			res = append(res, *k.(*Area))
		}
	}
	return res
}

type areaListResult struct {
	AreaList struct {
		WS2ListResponse
		Areas []struct {
			Area
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"area"`
	} `xml:"area-list"`
}
