package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aerogo/http/client"
)

type GSRes struct {
	URL      string
	Title    string
	Released string
}

// <li>
// <div class="img">
//   <a
// 	href="/category/fullmetal-alchemist-brotherhood-dub"
// 	title="Fullmetal Alchemist: Brotherhood (Dub)"
//   >
// 	<img
// 	  src="https://gogocdn.net/cover/fullmetal-alchemist-brotherhood-dub.png"
// 	  alt="Fullmetal Alchemist: Brotherhood (Dub)"
// 	/>
//   </a>
// </div>
// <p class="name">
//   <a
// 	href="/category/fullmetal-alchemist-brotherhood-dub"
// 	title="Fullmetal Alchemist: Brotherhood (Dub)"
// 	>Fullmetal Alchemist: Brotherhood (Dub)</a
//   >
// </p>
// <p class="released">Released: 2009</p>
// </li>

func searchGogo(searchterm string) string {
	searchquery := url.QueryEscape(searchterm)
	response, err := client.Get(BASE_URL + "/search.html?keyword=" + searchquery).End()
	if err != nil {
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.String()))
	if err != nil {
		return ""
	}

	var Res []GSRes

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		var eachres GSRes
		psel := s.Find(`p[class="name"]`)
		linksel := psel.Find(`a[href][title]`)
		eachres.URL = linksel.AttrOr("href", "")
		eachres.Title = linksel.AttrOr("title", "")
		relsel := s.Find(`p[class="released"]`)
		eachres.Released = strings.Replace(relsel.Text(), "Released: ", "", -1)
		if eachres.URL != "" && eachres.Title != "" {
			Res = append(Res, eachres)
		}
	})

	for _, eachres := range Res {
		fmt.Println(eachres.Title, eachres.Released, eachres.URL)
	}

	return ""
}
