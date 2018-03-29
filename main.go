package main

import (
	"docnota/Route"
	"net/http"
	"github.com/rs/cors"
	"docnota/Utils"
	"log"
	"time"
	"flag"
)

func main() {
	config := Utils.MainConfig
	var dir string

	r := Route.Route()

	flag.StringVar(&dir, "dir", ".", config.ServerConf.StaticDir)
	flag.Parse()

	http.Handle("/", r)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(dir))))
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
		AllowedHeaders: []string{"Access-Control-Allow-Origin", "Authorization", "Access-Control-Expose-Headers"},
		ExposedHeaders: []string{"Access-Control-Allow-Origin", "Authorization", "Access-Control-Expose-Headers"},

	})
	handler := c.Handler(r)
	server := &http.Server{
		Addr: config.ServerConf.Port,
		ReadTimeout: config.ServerConf.ReadTimeout * time.Second,
		WriteTimeout: config.ServerConf.WriteTimeout * time.Second,
		Handler:	handler,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}