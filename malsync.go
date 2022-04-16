package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type GogoAnime struct {
	Sub string
	Dub string
}

func gogoframanilist(iurl string) (ret GogoAnime, err error) {
	aid, err := getIDfromanilisturl(iurl)
	if err != nil {
		return ret, err
	}
	ret, err = getGogoAnimeLinks(aid, "anilist")
	if err != nil {
		return ret, err
	}
	return ret, nil
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

func getGogoAnimeLinks(id string, al string) (ret GogoAnime, err error) {
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
				ret.Dub = ""
			}
			ret.Dub = strings.TrimPrefix(u.Path, "/category/")
		} else {
			u, err := url.Parse(eachresp.URL)
			if err != nil {
				ret.Sub = ""
			}
			ret.Sub = strings.TrimPrefix(u.Path, "/category/")
		}
	}

	return
}
