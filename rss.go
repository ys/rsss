package main

import(
  "net/http"
  "fmt"
  "os"
  "strings"
  "encoding/xml"
)

type Rss struct {
  Channel Channel `xml:"channel"`
}

type Channel struct {
  Title         string `xml:"title"`
  Link          string `xml:"link"`
  Description   string `xml:"description"`
  Language      string `xml:"language"`
  LastBuildDate Date   `xml:"lastBuildDate"`
  Item1          []Item `xml:"item"`
  Item2          []Item `xml:"entry"`
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
  Updated     Date          `xml:"updated"`
  GUID        string        `xml:"guid"`
  ID        string          `xml:"id"`
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
  fmt.Println(resp.Header)
  defer resp.Body.Close()
  xmlDecoder := xml.NewDecoder(resp.Body)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  header := resp.Header.Get("Content-Type")
  var rss Rss
  var channel Channel
  if strings.Contains(header, "application/rss+xml") || strings.Contains(header, "text/xml") {
    err = xmlDecoder.Decode(&rss)
  } else {
    err = xmlDecoder.Decode(&channel)
  }
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  if strings.Contains(header, "application/rss+xml") || strings.Contains(header, "text/xml") {
    return rss.Channel
  } else {
    return channel
  }
}

func (i Item) PubDateFormatted() string {
  if i.PubDate != "" {
    return toNiceDate(i.PubDate)
  } else {
    return toNiceDate(i.Updated)
  }
}

func (i Channel) Item() []Item {
  if i.Item1 != nil {
    return i.Item1
  } else {
    return i.Item2
  }
}
