package main

import (
  "fmt"
)

func RunLogging(logs chan string) {
  for log := range logs {
    fmt.Println(log)
  }
}
