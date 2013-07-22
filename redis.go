package main

import(
  "os"
  "net/url"
  "github.com/garyburd/redigo/redis"
)

func RedisClient() redis.Conn {
  redisUrl := os.Getenv("REDIS_URL")
  if redisUrl == "" {
    redisUrl = "redis://127.0.0.1:6379"
  }
  redisUrl, pw := parseRedisUrl(redisUrl)
  c, err := redis.Dial("tcp", redisUrl)
  if err != nil {
    panic(err)
  }
   _, authErr := c.Do("AUTH", pw)
   if authErr != nil {
     // handle error
   }
  return c
}

func parseRedisUrl(redisUrl string) (string, string) {
  u, err := url.Parse(redisUrl)
  if err != nil {
    // handle error
  }
  pw := ""
  if u.User != nil {
    pw, _ = u.User.Password()
  }
  return u.Host, pw

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

