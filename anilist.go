package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/aerogo/http/client"
	"github.com/tidwall/gjson"
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
	ret, err = getGogoAnimeLinks(aid)
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

func getGogoAnimeLinks(id string) (ret GogoAnime, err error) {
	response, err := client.Get(fmt.Sprintf("https://raw.githubusercontent.com/MALSync/MAL-Sync-Backup/master/data/anilist/anime/%s.json", id)).End()
	if err != nil {
		log.Fatal(err)

	}
	result := gjson.Get(response.String(), "Pages.Gogoanime")
	subdone := false
	dubdone := false
	result.ForEach(func(key, value gjson.Result) bool {
		if strings.HasSuffix(key.String(), "dub") {
			val := value.Get("url")
			if !value.Exists() {
				println("no url")
			} else {
				u, err := url.Parse(val.String())
				if err != nil {
					ret.Dub = ""
				}
				ret.Dub = strings.TrimPrefix(u.Path, "/category/")
			}
			dubdone = true
		} else {
			val := value.Get("url")
			if !value.Exists() {
				println("no url")
			} else {
				u, err := url.Parse(val.String())
				if err != nil {
					ret.Dub = ""
				}
				ret.Sub = strings.TrimPrefix(u.Path, "/category/")
			}
			subdone = true
		}
		if subdone && dubdone {
			return false
		}
		return true // keep iterating

	})

	return ret, nil

}
