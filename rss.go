package main

import(
  "net/http"
  "fmt"
  "os"
  "encoding/xml"
  "github.com/garyburd/redigo/redis"
  "github.com/gorilla/mux"
  "encoding/json"
)

type RSS struct {
  Channel Channel `xml:"channel"`
}

type Channel struct {
  Title         string `xml:"title"`
  Link          string `xml:"link"`
  Description   string `xml:"description"`
  Language      string `xml:"language"`
  LastBuildDate Date   `xml:"lastBuildDate"`
  Item          []Item `xml:"item"`
}

type ItemEnclosure struct {
  URL  string `xml:"url,attr"`
  Type string `xml:"type,attr"`
}

type Item struct {
  Title       string        `xml:"title"`
  Link        string        `xml:"link"`
  Comments    string        `xml:"comments"`
  PubDate     Date          `xml:"pubDate"`
  GUID        string        `xml:"guid"`
  Category    []string      `xml:"category"`
  Enclosure   ItemEnclosure `xml:"enclosure"`
  Description string        `xml:"description"`
  Content     string        `xml:"content"`
}
type Date string

func AddRSSFeed(name string, url string) (interface{}, error) {
  c, err := redis.Dial("tcp", ":6379")
  defer c.Close()
  if err != nil {
    fmt.Println(err)
  }
  c.Send("HSET", "RSSS", name, url)
  c.Flush()
  v, err := c.Receive()

  return v, err
}

func getAllRSSS() (map[string]string) {
  c, err := redis.Dial("tcp", ":6379")
  defer c.Close()
  if err != nil {
    fmt.Println(err)
  }
  values, err := redis.Values(c.Do("HGETALL", "RSSS"))
  rsss := make(map[string]string)
  for i := 0; i < len(values); i += 2 {
    rsss[fmt.Sprintf("%s", values[i])] = fmt.Sprintf("%s", values[i+1])
  }
  return rsss
}

func AddItem(item Item, collection string) (interface{}, error) {
  fmt.Println(collection)
  c, err := redis.Dial("tcp", ":6379")
  defer c.Close()
  if err != nil {
    fmt.Println(err)
  }
  c.Send("LPUSH", "items", item)
  c.Flush()
  v, err := c.Receive()
  return v, err
}

func GetDistantRss(url string) Channel {
  resp, err := http.Get(os.Args[1])
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  defer resp.Body.Close()
  xmlDecoder := xml.NewDecoder(resp.Body)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  var rss RSS
  err = xmlDecoder.Decode(&rss)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  return rss.Channel
}

type Response map[string]string

func (r Response) String() (s string) {
  b, err := json.Marshal(r)
  if err != nil {
    s = ""
    return
  }
  s = string(b)
  return
}
func getRSSSHandle(w http.ResponseWriter, r *http.Request) {
  v := getAllRSSS()
  fmt.Fprint(w, Response(v))
  return
}

func setRSSHandle(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  r.ParseForm()
  fmt.Println(r.Form)

  name := r.Form["name"][0]
  url := r.Form["url"][0]
  v, err := AddRSSFeed(name, url)
  if err != nil {
    fmt.Fprint(w, "{'oops': true}")
    return
  }
  fmt.Println(v)
  resp := make(map[string]string)
  resp[name] = url
  fmt.Fprint(w, Response(resp))
  return
}

func updateRSSHandle(w http.ResponseWriter, r *http.Request) {
  fmt.Println("updateTRSSS")
  vars := mux.Vars(r)
  name := vars["name"]
  rss := GetDistantRss(name)
  for _, item := range rss.Item {
    fmt.Println(item.GUID)
    AddItem(item, name)
  }
}

func main() {
  r := mux.NewRouter()
  r.HandleFunc("/rsss", getRSSSHandle).Methods("GET")
  r.HandleFunc("/rsss", setRSSHandle).Methods("POST")
  r.HandleFunc("/rsss/{name}", updateRSSHandle).Methods("POST")
  http.Handle("/", r)
  http.ListenAndServe(":8080", nil)
}
