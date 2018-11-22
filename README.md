# zeitx
utilities to quickly bootstrap and http all-json service

## Usage

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

There are also json responses for errors as well as a handy querystring parameter `pretty` in order to easily inspect and debug:

```
-$ curl -i http://localhost:8090    
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Thu, 22 Nov 2018 13:37:24 GMT
Content-Length: 19

{"Version":"3.0.0"}%
```

```
-$ curl -i http://localhost:8090\?pretty
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
X-Pretty-Print: 1
Date: Thu, 22 Nov 2018 13:37:28 GMT
Content-Length: 26

{
    "Version": "3.0.0"
}%                                                                        
```


See main.go for a working example.
