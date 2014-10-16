package gomusicbrainz

// Tag is the common type for Tags.
type Tag struct {
	Count int    `xml:"count,attr"`
	Name  string `xml:"name"`
}

// TagResponse is the response type returned by tag request methods.
type TagResponse struct {
	WS2ListResponse
	Tags   []Tag
	Scores ScoreMap
}

// ResultsWithScore returns a slice of Tags with a specific score.
func (r *TagResponse) ResultsWithScore(score int) []Tag {
	var res []Tag
	for k, v := range r.Scores {
		if v == score {
			res = append(res, *k.(*Tag))
		}
	}
	return res
}

type tagListResult struct {
	TagList struct {
		WS2ListResponse
		Tags []struct {
			Tag
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"tag"`
	} `xml:"tag-list"`
}
