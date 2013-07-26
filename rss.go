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

type Link struct {
  Value string `xml:",chardata"`
  Attr  string `xml:"href,attr"`
}

type Item struct {
  Title       string        `xml:"title"`
  ILink       Link          `xml:"link"`
  Comments    string        `xml:"comments"`
  PubDate     Date          `xml:"pubDate"`
  Updated     Date          `xml:"updated"`
  GUID        string        `xml:"guid"`
  ID          string        `xml:"id"`
  Category    []string      `xml:"category"`
  Enclosure   ItemEnclosure `xml:"enclosure"`
  Description string        `xml:"description"`
  Content     string        `xml:"content"`
  Host        string
  HostLink    string
}

type Date string

func get(url string) *http.Response {
  resp, err := http.Get(url)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  return resp
}

func fetchRssFeed(url string) Channel {
  resp := get(url)
  defer resp.Body.Close()

  xmlDecoder := xml.NewDecoder(resp.Body)

  var rss Rss
  var channel Channel
  var err error
  isRss := responseIsRss(resp)
  if isRss {
    err = xmlDecoder.Decode(&rss)
  } else {
    err = xmlDecoder.Decode(&channel)
  }
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  if isRss {
    return rss.Channel
  } else {
    return channel
  }
}

func responseIsRss(response *http.Response) bool {
  header := response.Header.Get("Content-Type")
  return strings.Contains(header, "application/rss+xml") || strings.Contains(header, "text/xml")
}

func (i Item) PubDateFormatted() string {
  if i.PubDate != "" {
    return toNiceDate(i.PubDate)
  } else {
    return toNiceDate(i.Updated)
  }
}

func (i Item) Link() string {
  if i.ILink.Attr != "" {
    return i.ILink.Attr
  } else {
    return i.ILink.Value
  }
}

func (i Channel) Item() []Item {
  if i.Item1 != nil {
    return i.Item1
  } else {
    return i.Item2
  }
}
