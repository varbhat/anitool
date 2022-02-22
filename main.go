package main

import (
	"fmt"
)

var (
	BASE_URL = "https://gogoanime.fi"
)

func main() {
	//BASE_URL = finalURL(BASE_URL)
	//fmt.Println(BASE_URL)

	var iurl string
	fmt.Scanln(&iurl)

	// an, _ := gogoframanilist(iurl)
	// fmt.Println("Sub = ", an.Sub)
	// fmt.Println("Dub = ", an.Dub)
	// ep, _ := getEpsfromGogoID(an.Dub)
	// fmt.Println(ep)
	//fmt.Println(getDpageLink(iurl, "100"))

	decryptDLink("https://gogoplay.io/embedplus?id=MjgwNzY=&token=n-SJr3LziDVYK-tNM2fWHw&expires=1645533075")

}
