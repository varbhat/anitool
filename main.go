package main

import "fmt"

func main() {
	var abc string
	fmt.Scanln(&abc)

	// mon := gogoStreamLinks("https://gogoanime.fi", "midori-days", "5")
	// for each := range mon {
	// 	fmt.Println(each)
	// }

	// mon2 := Fplayer(abc)
	// for _, each := range mon2 {
	// 	fmt.Println(each.File, each.Label, each.Referer, each.Type)

	// }

	// mon3 := StreamSB(abc)
	// fmt.Println(mon3)

	mon4 := GoGoCDN(abc)
	for _, each := range mon4 {
		fmt.Println(each.File, each.Label, each.Referer, each.Type)

	}
}
