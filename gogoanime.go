package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type GogoAnime struct {
	Sub string
	Dub string
}

func (g *GogoAnime) FromAnilist(iurl string) (err error) {
	aid, err := getIDfromanilisturl(iurl)
	if err != nil {
		return err
	}
	g.Sub, g.Dub, err = getGogoAnimeLinks(aid, "anilist")
	if err != nil {
		return err
	}
	return nil
}

func (g *GogoAnime) FromMAL(iurl string) (err error) {
	aid, err := getIDfromMALurl(iurl)
	if err != nil {
		return err
	}
	g.Sub, g.Dub, err = getGogoAnimeLinks(aid, "myanimelist")
	if err != nil {
		return err
	}
	return nil
}

func (g *GogoAnime) GetLinks(aid string, epno string) (Ret chan Link) {
	var gogobaseurl string = "https://gogoanime.fi"
	gogobaseurl, err := finalURL(gogobaseurl)
	if err != nil {
		return
	}
	Ret = make(chan Link)
	go func() {
		var wg sync.WaitGroup
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
						eachlinku = "https:" + eachlinku
					}
					if strings.Contains(eachlinku, "gogo") || strings.Contains(eachlinku, "goload") {
						wg.Add(1)
						go func() {
							defer wg.Done()
							for _, eachretlink := range GoGoCDN(eachlinku) {
								Ret <- eachretlink
							}
						}()
					} else if strings.Contains(eachlinku, "sbplay") {
						wg.Add(1)
						go func() {
							defer wg.Done()
							for _, eachretlink := range StreamSB(eachlinku) {
								Ret <- eachretlink
							}
						}()
					} else if strings.Contains(eachlinku, "fplayer") || strings.Contains(eachlinku, "fembed") {
						wg.Add(1)
						go func() {
							defer wg.Done()
							for _, eachretlink := range Fplayer(eachlinku) {
								Ret <- eachretlink
							}
						}()
					}
				}
			} else {
				break
			}

			wg.Wait()
			close(Ret)
		}
	}()
	return
}

func getIDfromanilisturl(iurl string) (id string, err error) {
	u, err := url.Parse(iurl)
	if err != nil {
		return "", err
	}
	if u.Hostname() != "anilist.co" {
		return "", errors.New("invalid url")
	}

	params := strings.Split(u.Path, "/")
	paramlen := len(params)
	if paramlen < 2 {
		return "", errors.New("invalid url")

	}
	if params[1] != "anime" {
		return "", errors.New("invalid url")
	}
	return params[2], nil
}

func getIDfromMALurl(iurl string) (id string, err error) {
	u, err := url.Parse(iurl)
	if err != nil {
		return "", err
	}
	if u.Hostname() != "myanimelist.net" {
		return "", errors.New("invalid url")
	}

	params := strings.Split(u.Path, "/")
	paramlen := len(params)
	if paramlen < 2 {
		return "", errors.New("invalid url")

	}
	if params[1] != "anime" {
		return "", errors.New("invalid url")
	}
	return params[2], nil
}

func getGogoAnimeLinks(id string, al string) (sub string, dub string, err error) {
	var MSBResp struct {
		Pages struct {
			GGA map[string]struct {
				URL string `json:"url"`
			} `json:"Gogoanime"`
		} `json:"pages"`
	}
	response, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/MALSync/MAL-Sync-Backup/master/data/%s/anime/%s.json", al, id))
	if err != nil {
		fmt.Println(err)
	}
	respb, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(respb, &MSBResp)
	if err != nil {
		fmt.Println(err)
		return
	}
	for eachtitle, eachresp := range MSBResp.Pages.GGA {
		if strings.HasSuffix(eachtitle, "dub") {
			u, err := url.Parse(eachresp.URL)
			if err != nil {
				dub = ""
			}
			dub = strings.TrimPrefix(u.Path, "/category/")
		} else {
			u, err := url.Parse(eachresp.URL)
			if err != nil {
				sub = ""
			}
			sub = strings.TrimPrefix(u.Path, "/category/")
		}
	}

	return
}
