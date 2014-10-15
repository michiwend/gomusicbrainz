package gomusicbrainz

// ReleaseGroup groups several different releases into a single logical entity.
// Every release belongs to one, and only one release group. More informations
// at https://musicbrainz.org/doc/Release_Group
type ReleaseGroup struct {
	ID           string       `xml:"id,attr"`
	Type         string       `xml:"type,attr"`
	PrimaryType  string       `xml:"primary-type"`
	Title        string       `xml:"title"`
	ArtistCredit ArtistCredit `xml:"artist-credit"`
	Releases     []Release    `xml:"release-list>release"` // FIXME if important unmarshal count,attr
	Tags         []Tag        `xml:"tag-list>tag"`
}

// ReleaseGroupResponse is the response type returned by release group request
// methods.
type ReleaseGroupResponse struct {
	WS2ListResponse
	ReleaseGroups []ReleaseGroup
	Scores        ScoreMap
}

type releaseGroupListResult struct {
	ReleaseGroupList struct {
		WS2ListResponse
		ReleaseGroups []struct {
			ReleaseGroup
			Score int `xml:"http://musicbrainz.org/ns/ext#-2.0 score,attr"`
		} `xml:"release-group"`
	} `xml:"release-group-list"`
}
