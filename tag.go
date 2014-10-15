package gomusicbrainz

// Tag is the common type for Tags.
type Tag struct {
	Count int    `xml:"count,attr"`
	Name  string `xml:"name"`
	Score int    `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
}

// TagResponse is the response type returned by tag request methods.
type TagResponse struct {
	WS2ListResponse
	Tags []Tag `xml:"tag"`
}
