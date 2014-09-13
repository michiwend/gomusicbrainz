gomusicbrainz
=============

a Go (Golang) MusicBrainz WS2 client library - work in progress. [![Build Status](https://travis-ci.org/michiwend/gomusicbrainz.svg?branch=master)](https://travis-ci.org/michiwend/gomusicbrainz)

![gopherbrainz Oo](misc/gopherbrainz.png)

## Current state
Currently only search requests are supported. Browse and lookup requests will
follow as soon as all search requests are covered.

## Example usage
This simple example requests all artist matching "bonobo" or "Parov Stelar"
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
## Full documentation
All search request follow the [Apache Lucene syntax](https://lucene.apache.org/core/4_3_0/queryparser/org/apache/lucene/queryparser/classic/package-summary.html#package_description). Please head over to the [MusicBrainz website] (https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search) for more information about all possible query-fields.

Documentation for this package can be found on godoc.org (badge).

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/michiwend/gomusicbrainz)
