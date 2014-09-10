package gomusicbrainz

import (
	"fmt"
	"testing"
)

func TestFindArtist(t *testing.T) {

	mbc := New("http://musicbrainz.org/ws/2/")

	artists, err := mbc.SearchArtist("kruder")

	if err != nil {
		t.Error(err)
	}

	for _, artist := range artists {
		fmt.Println(artist.Name)

	}

}
