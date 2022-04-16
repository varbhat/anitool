package main

import "fmt"

func main() {
	var abc string
	fmt.Scanln(&abc)

	mon := gogoStreamLinks("https://gogoanime.fi", "midori-days", "5")
	for each := range mon {
		fmt.Println(each)
	}
}
