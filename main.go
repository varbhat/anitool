package main

import "fmt"

func main() {
	var iurl string
	fmt.Scanln(&iurl)
	res := searchGogoAll("https://gogoanime.fi", iurl)
	// for _, each := range res {
	// 	fmt.Print(each.File + "  " + each.Label + "  " + each.Type + " " + each.Referer)
	// 	fmt.Println()
	// }
	for each := range res {
		fmt.Println(each.Released, each.Title, each.URL)
	}
	fmt.Println("THE End")

}
