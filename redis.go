package main

import(
  "os"
  "github.com/garyburd/redigo/redis"
)

func RedisClient() redis.Conn {
  url := os.Getenv("REDIS_URL")
  if url == "" {
    url = ":6379"
  }
  c, err := redis.Dial("tcp", url)
  if err != nil {
    panic(err)
  }
  return c
}

func redisDo(cmd string, args ...interface{}) interface{} {
  c := RedisClient()
  defer c.Close()
  value, err := c.Do(cmd, args...)
  if err != nil {
    panic(err)
  }
  return value
}

func toString(data interface{}) string {
  value, err := redis.String(data, nil)
  if err != nil {
    panic(err)
  }
  return value
}

func toStrings(data interface{}) []string {
  value, err := redis.Strings(data, nil)
  if err != nil {
    panic(err)
  }
  return value
}

func toValues(data interface{}) []interface{} {
  value, err := redis.Values(data, nil)
  if err != nil {
    panic(err)
  }
  return value
}

