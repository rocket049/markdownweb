package main

import (
	"log"
	"net/http"
	"path"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a := strings.Split(r.Host, ":")
		url1 := path.Join("https://"+a[0], r.URL.Path)
		http.Redirect(w, r, url1, 301)
		//log.Println(url1)
	})
	err := http.ListenAndServe("127.0.0.1:8090", nil)
	log.Println(err)
}
