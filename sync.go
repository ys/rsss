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
      importRss(url)
    }
    time.Sleep(6000 * 1000 * time.Millisecond)
  }
}
