package main

import(
  "net/http"
)

func main() {
  logs := make(chan string, 10000)
  go RunLogging(logs)

  go SyncRssFeeds()

  handler := routerHandlerFunc(router())
  handler = wrapLogging(handler, logs)
  logs <- "Listening on port 8888"
  http.ListenAndServe(":8888", handler)
}
