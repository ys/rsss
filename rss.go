package main

import(
  "net/http"
  "fmt"
  "os"
  "encoding/xml"
)

type RSS struct {
  Channel Channel `xml:"channel"`
}

type Channel struct {
  Title         string `xml:"title"`
  Link          string `xml:"link"`
  Description   string `xml:"description"`
  Language      string `xml:"language"`
  LastBuildDate Date   `xml:"lastBuildDate"`
  Item          []Item `xml:"item"`
}

type ItemEnclosure struct {
  URL  string `xml:"url,attr"`
  Type string `xml:"type,attr"`
}

type Item struct {
  Title       string        `xml:"title"`
  Link        string        `xml:"link"`
  Comments    string        `xml:"comments"`
  PubDate     Date          `xml:"pubDate"`
  GUID        string        `xml:"guid"`
  Category    []string      `xml:"category"`
  Enclosure   ItemEnclosure `xml:"enclosure"`
  Description string        `xml:"description"`
  Content     string        `xml:"content"`
  Host        string
}

type Date string

func fetchRssFeed(url string) Channel {
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer resp.Body.Close()
  xmlDecoder := xml.NewDecoder(resp.Body)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  var rss RSS
  err = xmlDecoder.Decode(&rss)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  return rss.Channel
}

func (i Item) PubDateFormatted() string {
  return toNiceDate(i.PubDate)
}
