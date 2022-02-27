package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/aerogo/http/client"
)

func gogoStreamLinks(gogobaseurl string, aid string, epno string) (ret chan string) {
	ret = make(chan string)
	defer close(ret)
	paramlist := []string{"-episode-%s", "-%s", "-episode-%s-1", "-camrip-episode-%s"}
	for _, eachparam := range paramlist {
		response, err := http.Get(fmt.Sprintf(gogobaseurl+"/"+aid+eachparam, epno))
		if err != nil {
			continue
		}
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			continue
		}

		if doc.Find(".entry-title").Text() != "404" {
			dv := doc.Find("a[data-video]")
			for _, eachlink := range dv.Nodes {
				linko := goquery.NewDocumentFromNode(eachlink)
				eachlinku := linko.AttrOr("data-video", "")
				if strings.HasPrefix(eachlinku, "//") {
					ret <- "https:" + eachlinku
				} else {
					ret <- eachlinku
				}
			}
		} else {
			break
		}

	}
	return
}

type GSRes struct {
	URL      string
	Title    string
	Released string
}

func getEpsCount(gogobaseurl string, gid string) (ret int, err error) {
	ret = 0

	response, err := client.Get(gogobaseurl + "/category/" + gid).End()
	if err != nil {
		return 0, err

	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.String()))
	if err != nil {
		return 0, err
	}

	avb := doc.Find(".anime_video_body")
	for _, eachnode := range avb.Nodes {
		gs := goquery.NewDocumentFromNode(eachnode)
		epnodes := gs.Find("a")
		for _, eachepn := range epnodes.Nodes {
			epngs := goquery.NewDocumentFromNode(eachepn)
			ep := epngs.AttrOr("ep_end", "0")
			episode, err := strconv.Atoi(ep)
			if err == nil {
				if episode > ret {
					ret = episode
				}

			}

		}
	}
	return ret, nil
}

func getPagecount(gogobaseurl string, searchterm string) (pagecount int) {
	searchquery := url.QueryEscape(searchterm)
	response, err := http.Get(gogobaseurl + "/search.html?keyword=" + searchquery)
	if err != nil {
		return 0
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
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

func getGogoSearchRes(gogobaseurl string, searchterm string, page int) (Res chan GSRes) {
	Res = make(chan GSRes)
	go func() {
		searchquery := url.QueryEscape(searchterm)
		response, err := client.Get(gogobaseurl + "/search.html?keyword=" + searchquery + "&page=" + strconv.Itoa(page)).End()
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
				Res <- eachres
			}
		})

		close(Res)
	}()

	return Res
}

func searchGogoAll(gogobaseurl string, searchterm string) (Res chan GSRes) {
	Res = make(chan GSRes)

	go func() {

		pagecount := getPagecount(gogobaseurl, searchterm)
		var wg sync.WaitGroup
		for ep := 0; ep <= pagecount; ep++ {
			retch := getGogoSearchRes(gogobaseurl, searchterm, pagecount)
			wg.Add(1)
			go func(c <-chan GSRes) {
				defer wg.Done()
				for v := range c {
					Res <- v
				}

			}(retch)
		}
		wg.Wait()
		close(Res)
	}()
	return Res

}
