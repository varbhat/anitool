package main

import "fmt"

func main() {
	var ga GogoAnime
	mon := ga.GetLinks("midori-days", "5")
	var Ret []Link
	for each := range mon {
		Ret = append(Ret, each)
		fmt.Println(each.File, each.Label, each.Referer, each.Type)
	}
	var inp int
	fmt.Scanln(&inp)
	Ret[inp].Play()
}
