package main

import(
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/hoisie/mustache"
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

type statusCapturingResponseWriter struct {
  status int
  http.ResponseWriter
}

func (w statusCapturingResponseWriter) WriteHeader(s int) {
  w.status = s
  w.ResponseWriter.WriteHeader(s)
}

func routerHandlerFunc(router *mux.Router) http.HandlerFunc {
  return func(res http.ResponseWriter, req *http.Request) {
    router.ServeHTTP(res, req)
  }
}

func getRssHandle(w http.ResponseWriter, r *http.Request) {
  v := getAllRssFeeds()
  fmt.Fprint(w, Response(v))
  return
}

func setRssHandle(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  r.ParseForm()
  name := r.Form.Get("name")
  url := r.Form.Get("url")
  AddRssFeed(name, url)
  resp := make(map[string]string)
  resp[name] = url
  fmt.Fprint(w, Response(resp))
  return
}

func getRssItemsHandle(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  items := GetAllItems()
  resp := make(map[string][]Item)
  resp["items"] = items
  fmt.Fprint(w, ResponseArray(resp))
  return
}

func getRssHtmlHandle(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  items := GetAllItems()
  resp := make(map[string][]Item)
  resp["items"] = items
  fmt.Fprint(w, mustache.RenderFileInLayout("index.html.mustache", "layout.html.mustache", resp))
}

func newRssHtmlHandle(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  fmt.Fprint(w, mustache.RenderFileInLayout("new.html.mustache", "layout.html.mustache", nil))
}

func router() *mux.Router {
  r := mux.NewRouter()
  r.HandleFunc("/rsss", getRssHtmlHandle).Methods("GET")
  r.HandleFunc("/rsss/feeds", getRssHandle).Methods("GET")
  r.HandleFunc("/rsss/all", getRssItemsHandle).Methods("GET")
  r.HandleFunc("/rsss/new", newRssHtmlHandle).Methods("GET")
  r.HandleFunc("/rsss", setRssHandle).Methods("POST")
  r.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
  return r
}
