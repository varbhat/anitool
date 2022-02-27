package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/aerogo/http/client"
)

type FplayerResp struct {
	Data []Link `json:"data"`
}

func Fplayer(iurl string) []Link {
	apiurl := strings.Replace(iurl, "/v/", "/api/source/", -1)
	response, err := client.Post(apiurl).End()
	if err != nil {
		log.Fatal(err)

	}
	var Fresp FplayerResp
	if err := json.Unmarshal(response.Bytes(), &Fresp); err != nil {
		return []Link{}
	}

	for _, eachSource := range Fresp.Data {
		eachSource.Referer = iurl
	}

	return Fresp.Data
}
