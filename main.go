package main

import (
	"flag"
	"laszlo/zeitx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	ver        = "3.0.0"
	configFile = flag.String("config", "./config.yml", "config filename")
)

func main() {
	flag.Parse()
	cfg, err := zeitx.NewConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	srv := zeitx.NewHTTPServer(*cfg)
	srv.GETRoute("/", index)
	srv.ListenAndServe()

	sc := make(chan os.Signal)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	<-sc
}

func index(w http.ResponseWriter, r *http.Request) {
	zeitx.OkJSON(w, r, &struct{ Version string }{Version: ver})
}
