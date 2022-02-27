package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/aerogo/http/client"
)

func streamSB(iurl string) string {
	ssburl := "https://sbplay2.com"
	jsonlink := ssburl + "/sources40/7361696b6f757c7c%s7c7c7361696b6f757c7c73747265616d7362/7361696b6f757c7c363136653639366436343663363136653639366436343663376337633631366536393664363436633631366536393664363436633763376336313665363936643634366336313665363936643634366337633763373337343732363536313664373336327c7c7361696b6f757c7c73747265616d7362"

	u, err := url.Parse(iurl)
	if err != nil {
		return ""
	}
	params := strings.Split(u.Path, "/")
	paramlen := len(params)
	if paramlen < 2 {
		return ""

	}
	if params[1] != "e" {
		return ""
	}
	sourceid := hex.EncodeToString([]byte(params[2]))
	jsonApi := fmt.Sprintf(jsonlink, sourceid)
	response, err := client.Get(jsonApi).Header("watchsb", "streamsb").End()
	if err != nil {
		log.Fatal(err)

	}
	fmt.Println(response.String())
	return ""
}
