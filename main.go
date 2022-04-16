package main

import "fmt"

func main() {
	// mon := gogoStreamLinks("midori-days", "5")
	// var Ret []Link
	// for each := range mon {
	// 	Ret = append(Ret, each)
	// 	fmt.Println(each.File, each.Label, each.Referer, each.Type)
	// }
	// var inp int
	// fmt.Scanln(&inp)
	// Ret[inp].Play()

	var ga GogoAnime
	ga.FromAnilist("https://anilist.co/anime/20767/Date-A-Live-II-Kurumi-Star-Festival/")
	fmt.Println(ga.Dub, ga.Sub)
}
