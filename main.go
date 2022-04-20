package main

import (
	"flag"
	"fmt"
)

func main() {

	var animename string
	var episode string

	flag.StringVar(&animename, "anime", "", "Anime Name")
	flag.StringVar(&episode, "episode", "", "Episode")
	flag.Parse()
	if animename == "" {
		fmt.Println("Enter anime name")
		fmt.Scanln(&animename)
	}
	if episode == "" {
		fmt.Println("Enter Episode")
		fmt.Scanln(&episode)
	}
	var ga GogoAnime
	mon := ga.GetLinks(animename, episode)
	var Ret []Link
	count := 0
	for each := range mon {
		Ret = append(Ret, each)
		fmt.Println(count, each.File, each.Label, each.Referer, each.Type)
		count++
		fmt.Println()
	}
	var inp int
	fmt.Scanln(&inp)
	Ret[inp].Play()
}
