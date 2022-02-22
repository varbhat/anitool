package main

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aerogo/http/client"
)

func getEpsfromGogoID(gid string) (ret int, err error) {
	ret = 0

	response, err := client.Get(BASE_URL + "/category/" + gid).End()
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
