# zeitx
utilities to quickly bootstrap and http all-json service


```
func main() {
  ...
  
  cfg, _ := zeitx.NewConfig(*configFile)
  srv := zeitx.NewHTTPServer(*cfg)
  srv.GETRoute("/", index)
  srv.ListenAndServe()
  
  ...
}

func index(w http.ResponseWriter, r *http.Request) {
  zeitx.OkJSON(w, r, &struct{ Version string }{Version: ver})
}
```

See main.go for a working example.
