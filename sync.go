package main

import(
  "fmt"
  "time"
)
func SyncRssFeeds() {
  for {
    fmt.Println(time.Now())
    rssFeeds := getRssFeedUrls()
    fmt.Println(rssFeeds)
    for _, url := range rssFeeds {
      go importRss(url)
    }
    time.Sleep(6000 * 1000 * time.Millisecond)
  }
}
