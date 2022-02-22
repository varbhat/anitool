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

	an, _ := gogoframanilist(iurl)
	fmt.Println("Sub = ", an.Sub)
	fmt.Println("Dub = ", an.Dub)

}
