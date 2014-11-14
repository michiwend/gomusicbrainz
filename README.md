# gomusicbrainz [![License MIT](http://img.shields.io/badge/License-MIT-lightgrey.svg?style=flat-square)](http://opensource.org/licenses/MIT) [![GoDoc](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/michiwend/gomusicbrainz) [![GoWalker](http://img.shields.io/badge/api-GoWalker-green.svg?style=flat-square)](https://gowalker.org/github.com/michiwend/gomusicbrainz) [![Build Status](http://img.shields.io/travis/michiwend/gomusicbrainz.svg?style=flat-square)](https://travis-ci.org/michiwend/gomusicbrainz) 

a Go (Golang) MusicBrainz WS2 client library - a work in progress.

![gopherbrainz Oo](https://raw.githubusercontent.com/michiwend/gomusicbrainz/master/misc/gopherbrainz.png)

## Current state
Currently GoMusicBrainz provides methods to perform search and lookup requests. Browse requests are not supported yet.

## Installation
```bash
$ go get github.com/michiwend/gomusicbrainz
```

## Search Requests
GoMusicBrainz provides a search method for every WS2 search request in the form:
```Go
func (*WS2Client) Search<ENTITY>(searchTerm, limit, offset) (<ENTITY>SearchResponse, error)
```
searchTerm follows the Apache Lucene syntax and can either contain multiple
fields with logical operators or just a simple search string. Please refer to
[lucene.apache.org](https://lucene.apache.org/core/4_3_0/queryparser/org/apache/lucene/queryparser/classic/package-summary.html#package_description)
for more details on the lucene syntax. In addition the
[MusicBrainz website] (https://musicbrainz.org/doc/Development/XML_Web_Service/Version_2/Search)
provides information about all possible query-fields.

### Example
This example demonstrates a simple search requests to find the artist
*Parov Stelar*. You can find it as a runnable go program in the samples folder.
```Go
// create a new WS2Client.
client := gomusicbrainz.NewWS2Client(
    "https://musicbrainz.org/ws/2",
    "A GoMusicBrainz example",
    "0.0.1-beta",
    "http://github.com/michiwend/gomusicbrainz")

// Search for some artist(s)
resp, _ := client.SearchArtist(`artist:"Parov Stelar"`, -1, -1)

// Pretty print Name and score of each returned artist.
for _, artist := range resp.Artists {
    fmt.Printf("Name: %-25sScore: %d\n", artist.Name, resp.Scores[artist])
}
```
the above code will produce the following output:
```
Name: Parov Stelar             Score: 100
Name: Parov Stelar Trio        Score: 80
Name: Parov Stelar & the Band  Score: 70
```

## Lookup Requests
GoMusicBrainz provides two ways to perform lookup requests: Either the specific
lookup method that is implemented for each entity that has a lookup endpoint in
the form
```Go
func(*WS2Client) Lookup<ETITY>(id MBID) (*<ENTITY>, error)
```
or the common lookup method if you already have an entity (with MBID) that
implements the MBLookupEntity interface:
```Go
func(*MBLookupEntity) Lookup(entity MBLookupEntity, inc []string) error
```

### Example
The following example demonstrates the (specific) LookupArtist method. You can
find it as a runnable go program in the samples folder.

```Go
// create a new WS2Client.
client, _ := gomusicbrainz.NewWS2Client(
    "https://musicbrainz.org/ws/2",
    "A GoMusicBrainz example",
    "0.0.1-beta",
    "http://github.com/michiwend/gomusicbrainz")

// Lookup artist by id.
artist, err := client.LookupArtist("9a709693-b4f8-4da9-8cc1-038c911a61be", nil)

if err != nil {
    fmt.Println(err)
    return
}

fmt.Printf("%+v", artist)
```

## Package Documentation
Full documentation for this package can be found at
[GoDoc](https://godoc.org/github.com/michiwend/gomusicbrainz)
and  [GoWalker](https://gowalker.org/github.com/michiwend/gomusicbrainz)
