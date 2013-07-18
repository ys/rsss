package main

import(
  "fmt"
  "github.com/garyburd/redigo/redis"
  "encoding/json"
)

func RedisClient() redis.Conn {
  c, err := redis.Dial("tcp", ":6379")
  if err != nil {
    panic(err)
  }
  return c
}

func AddRSSFeed(name string, url string) (interface{}, error) {
  c := RedisClient()
  defer c.Close()
  c.Send("HSET", "RSSS", name, url)
  c.Flush()
  v, err := c.Receive()
  return v, err
}

func GetRSSUrl(name string) (string) {
  c := RedisClient()
  defer c.Close()
  url, _ := redis.String(c.Do("HGET", "RSSS", name))
  return url
}

func getAllRSSS() (map[string]string) {
  c:= RedisClient()
  defer c.Close()
  values, err := redis.Values(c.Do("HGETALL", "RSSS"))
  if err != nil {
    panic(err)
  }
  rsss := make(map[string]string)
  for i := 0; i < len(values); i += 2 {
    rsss[fmt.Sprintf("%s", values[i])] = fmt.Sprintf("%s", values[i+1])
  }
  return rsss
}

func GetAllItems() ([]Item) {
  c:= RedisClient()
  defer c.Close()
  values, err := redis.Strings(c.Do("LRANGE", "items", 0, 10))
  if err != nil {
    panic(err)
  }
  arr := make([]Item, 0)
  for _, value := range values {
    var item Item
    fmt.Println(value)
    json.Unmarshal([]byte(value), &item)
    arr = append(arr, item)
  }
  return arr
}

func AddItem(item Item) (interface{}, error) {
  c := RedisClient()
  defer c.Close()
  item_json, _ := json.Marshal(item)
  c.Send("LPUSH", "items", item_json)
  c.Flush()
  v, err := c.Receive()
  return v, err
}
