package main

import(
  "net/http"
)

func main() {
  logs := make(chan string, 10000)
  go RunLogging(logs)

  r := Router()
  http.Handle("/", r)
  logs <- "Listening on port 8888"
  http.ListenAndServe(":8888", nil)
}
