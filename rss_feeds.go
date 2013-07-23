package main

import(
  "fmt"
  "encoding/json"
)

func AddRssFeed(name string, url string) interface{} {
  v := redisDo("HSET", "RSSS", name, url)
  importRss(url)
  return v
}

func GetRssUrl(name string) (string) {
  return toString(redisDo("HGET", "RSSS", name))
}

func getRssFeedUrls() []string {
  return toStrings(redisDo("HVALS", "RSSS"))
}

func cleanItems() interface{} {
  return redisDo("DEL", "items")
}

func getAllRssFeeds() (map[string]string) {
  values := toStrings(redisDo("HGETALL", "RSSS"))
  rss_feeds := make(map[string]string)
  for i := 0; i < len(values); i += 2 {
    rss_feeds[values[i]] =  values[i+1]
  }
  return rss_feeds
}

func GetAllItems() ([]Item) {
  values := toStrings(redisDo("ZREVRANGE", "items", 0, -1))
  items := make([]Item, 0)
  for _, value := range values {
    var item Item
    json.Unmarshal([]byte(value), &item)
    items = append(items, item)
  }
  return items
}

func AddItem(item Item) interface{} {
  item_json, _ := json.Marshal(item)
  var score int64
  if item.PubDate == "" {
    fmt.Println(item.Host)
    score = toUnix(item.Updated)
  } else {
    score = toUnix(item.PubDate)
  }
  return redisDo("ZADD", "items", score , item_json)
}

func importRss(rssUrl string) {
  rss := fetchRssFeed(rssUrl)
  for _, item := range rss.Item() {
    item.Host = rss.Title
    AddItem(item)
  }
}
