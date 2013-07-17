package main

import(
  "net/http"
  "fmt"
  "os"
  "encoding/xml"
  "github.com/garyburd/redigo/redis"
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
}
type Date string

func AddItem(item Item) {
  c, err := redis.Dial("tcp", ":6379")
  defer c.Close()
  if err != nil {
    fmt.Println(err)
  }
  c.Send("LPUSH", "items", item)
  c.Flush()
  v, err := c.Receive()
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(v)
}

func GetRss(url string) Channel {
  resp, err := http.Get(os.Args[1])
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

func main() {
  rss := GetRss(os.Args[1])
  for _, item := range rss.Item {
    fmt.Println(item.Title)
    AddItem(item)
  }
}
