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

type tagListResult struct {
	TagList struct {
		WS2ListResponse
		Tags []struct {
			Tag
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"tag"`
	} `xml:"tag-list"`
}
