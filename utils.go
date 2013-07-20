package main

import (
  "time"
)

func toTime(timeString Date) time.Time {
  val, err := time.Parse("2006-01-02T15:04:05-07:00", string(timeString))
  if err != nil {
    val, _ = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", string(timeString))
  }
  return val
}
func toUnix(timeString Date) int64 {
  return toTime(timeString).Unix()
}

func toNiceDate(timeString Date) string {
  t := toTime(timeString)
  return t.Format(time.RFC1123)
}
