package main

import (
	"log"
	"net/http"
)

func finalURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("Couldn't HTTP Get URL: " + url)
	}
	return resp.Request.URL.String()
}
