package main

import(
  "encoding/json"
  "github.com/gorilla/mux"
  "net/http"
  "fmt"
)

type Response map[string]string
type ResponseArray map[string][]Item

func (r Response) String() (s string) {
  b, err := json.Marshal(r)
  if err != nil {
    s = ""
    return
  }
  s = string(b)
  return
}

func (r ResponseArray) String() (s string) {
  b, err := json.Marshal(r)
  if err != nil {
    s = ""
    return
  }
  s = string(b)
  return
}

func GetRSSHandle(w http.ResponseWriter, r *http.Request) {
  v := getAllRSSS()
  fmt.Fprint(w, Response(v))
  return
}

func SetRSSHandle(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  r.ParseForm()
  fmt.Println(r.Form)

  name := r.Form["name"][0]
  url := r.Form["url"][0]
  v, err := AddRSSFeed(name, url)
  if err != nil {
    fmt.Fprint(w, v)
    return
  }
  fmt.Println(v)
  resp := make(map[string]string)
  resp[name] = url
  fmt.Fprint(w, Response(resp))
  return
}

func UpdateRSSHandle(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  name := vars["name"]
  url := GetRSSUrl(name)
  fmt.Println(url)
  rss := GetRss(url)
  for _, item := range rss.Item {
    AddItem(item)
  }
}

func GetRSSItemsHandle(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  items := GetAllItems()
  resp := make(map[string][]Item)
  resp["items"] = items
  fmt.Fprint(w, ResponseArray(resp))
  return
}


func Router() *mux.Router {
  r := mux.NewRouter()
  r.HandleFunc("/rsss", GetRSSHandle).Methods("GET")
  r.HandleFunc("/rsss", SetRSSHandle).Methods("POST")
  r.HandleFunc("/rsss/all", GetRSSItemsHandle).Methods("GET")
  r.HandleFunc("/rsss/{name}", UpdateRSSHandle).Methods("POST")
  return r
}
