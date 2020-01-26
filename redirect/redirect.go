package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		a := strings.Split(r.Host, ":")
		url1 := "https://" + a[0]
		http.Redirect(w, r, url1, 301)
		//log.Println(url1)
	})
	err := http.ListenAndServe(":8080", nil)
	log.Println(err)
}
