
// convert.go: Convert rss.opml to YAML list of feed URLs under 'feeds' key
package main

import (
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "os"
    "gopkg.in/yaml.v3"
)

type OPML struct {
    Body OPMLBody `xml:"body"`
}

type OPMLBody struct {
    Outlines []Outline `xml:"outline"`
}

type Outline struct {
    XMLUrl  string    `xml:"xmlUrl,attr"`
    Outlines []Outline `xml:"outline"`
}

type Feed struct {
    URL string `yaml:"url"`
}

type Feeds struct {
    Feeds []Feed `yaml:"feeds"`
}

func extractFeeds(outlines []Outline, feeds *[]Feed) {
    for _, o := range outlines {
        if o.XMLUrl != "" {
            *feeds = append(*feeds, Feed{URL: o.XMLUrl})
        }
        if len(o.Outlines) > 0 {
            extractFeeds(o.Outlines, feeds)
        }
    }
}

func main() {
    // Read OPML file
    data, err := ioutil.ReadFile("rss.opml")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading rss.opml: %v\n", err)
        os.Exit(1)
    }

    // Parse OPML
    var opml OPML
    if err := xml.Unmarshal(data, &opml); err != nil {
        fmt.Fprintf(os.Stderr, "Error parsing OPML: %v\n", err)
        os.Exit(1)
    }

    // Extract feeds
    var feedsList []Feed
    extractFeeds(opml.Body.Outlines, &feedsList)

    // Marshal to YAML
    feeds := Feeds{Feeds: feedsList}
    yamlData, err := yaml.Marshal(&feeds)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error marshaling YAML: %v\n", err)
        os.Exit(1)
    }

    // Write to file
    if err := ioutil.WriteFile("feeds.yaml", yamlData, 0644); err != nil {
        fmt.Fprintf(os.Stderr, "Error writing feeds.yaml: %v\n", err)
        os.Exit(1)
    }
    fmt.Println("feeds.yaml created successfully.")
}
