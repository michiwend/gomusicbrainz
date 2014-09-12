gomusicbrainz
=============

a musicbrainz client library - work in progress

![gopherbrainz Oo](misc/gopherbrainz.png)

## Example usage
With SearchArtist you can query MusicBrainz' Search Server for some Artist.
searchTerm can be any string that follows [Apaches Lucene Syntax](https://lucene.apache.org/core/4_3_0/queryparser/org/apache/lucene/queryparser/classic/package-summary.html#package_description). This applies for every Search method.
```Go
import "github.com/michiwend/gomusicbrainz"

// create a new WS2 client
client := gomusicbrainz.NewWS2Client()
// Search for some artists
artists, _ := client.SearchArtist(`bonobo OR "Parov Stelar"`, -1, -1)

// Pretty print Name and Id of each returned artist.
for _, artist := range artists {
	fmt.Printf("Name: %-25s ID: %s\n", artist.Name, artist.Id)
}
```
