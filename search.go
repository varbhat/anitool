package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/aerogo/http/client"
)

type GSRes struct {
	URL      string
	Title    string
	Released string
}

func getPagecount(searchterm string) (pagecount int) {
	searchquery := url.QueryEscape(searchterm)
	response, err := client.Get(BASE_URL + "/search.html?keyword=" + searchquery).End()
	if err != nil {
		return 0
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.String()))
	if err != nil {
		return 0
	}
	doc.Find("a[data-page]").Each(func(i int, s *goquery.Selection) {
		pg, err := strconv.Atoi(s.Text())
		if err == nil {
			if pg >= pagecount {
				pagecount = pg
			}
		}
	})
	return
}

func getGogo(searchterm string, page int) (Res []GSRes) {
	searchquery := url.QueryEscape(searchterm)
	response, err := client.Get(BASE_URL + "/search.html?keyword=" + searchquery + "&page=" + strconv.Itoa(page)).End()
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.String()))
	if err != nil {
		return
	}

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

	return
}

func searchGogo(searchterm string) string {
	pagecount := getPagecount(searchterm)
	var wg sync.WaitGroup
	for ep := 0; ep <= pagecount; ep++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			getGogo(searchterm, pagecount)

		}()
	}
	wg.Wait()
	return ""
}
