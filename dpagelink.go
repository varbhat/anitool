package main

import (
	"bytes"
	"crypto/aes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aerogo/http/client"
)

type Link struct {
	File  string
	Label string
	Type  string
}

type Fresponse struct {
	Source   []Link `json:"source"`
	SourceBk []Link `json:"source_bk"`
}

func getDpageLink(aid string, epno string) (ret []string) {
	paramlist := []string{"-episode-%s", "-%s", "-episode-%s-1", "-camrip-episode-%s"}
	for _, eachparam := range paramlist {
		response, err := client.Get(fmt.Sprintf(BASE_URL+"/"+aid+eachparam, epno)).End()
		if err != nil {
			continue
		}
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.String()))
		if err != nil {
			continue
		}

		if doc.Find(".entry-title").Text() != "404" {
			dv := doc.Find("a[data-video]")
			for _, eachlink := range dv.Nodes {
				linko := goquery.NewDocumentFromNode(eachlink)
				// if linko.AttrOr("rel", "rel") != "100" {
				// 	continue
				// }
				eachlinku := linko.AttrOr("data-video", "")
				if strings.HasPrefix(eachlinku, "//") {
					ret = append(ret, "https:"+eachlinku)
				} else {
					ret = append(ret, eachlinku)
				}
			}
		} else {
			break
		}

	}
	return
}

func decryptDLink(iurl string) []Link {
	ajax_url := "https://gogoplay.io/encrypt-ajax.php"
	var rtime string
	var secret_key string
	var ajax string

	// Get video id
	iUrl, err := url.Parse(iurl)
	if err != nil {
		return []Link{}
	}
	video_id := iUrl.Query().Get("id")

	response, err := client.Get(iurl).End()
	if err != nil {
		return []Link{}
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.String()))
	if err != nil {
		return []Link{}
	}

	encrypted := doc.Find("script[data-name='crypto']").AttrOr("data-value", "")
	iv := doc.Find("script[data-name='ts']").AttrOr("data-value", "")
	if encrypted == "" || iv == "" {
		secret_key = "3235373436353338353932393338333936373634363632383739383333323838"
		iv = "34323036393133333738303038313335"
		rtime = "69420691337800813569"
		ajax = AES256Encrypt(secret_key, iv, video_id)
	} else {
		rtime = "00000000000000000000"
		secret_key = aes256decrypt(encrypted, []byte(iv+iv), []byte(iv), aes.BlockSize)
		ajax = aes256encrypt([]byte(video_id), []byte(secret_key), []byte("0000000000000000"), aes.BlockSize)
	}

	var rbody = []byte(fmt.Sprintf("id=%s&time=%s", ajax, rtime))
	req, _ := http.NewRequest("POST", ajax_url, bytes.NewBuffer(rbody))
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, _ := client.Do(req)
	//defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	//jsonResp := string(body)

	var Fresp Fresponse

	if err := json.Unmarshal(body, &Fresp); err != nil {
		return []Link{}
	}

	for _, eachSource := range Fresp.Source {
		fmt.Println("File ", eachSource.File)
		fmt.Println("Label ", eachSource.Label)
		fmt.Println("Type ", eachSource.Type)
	}

	for _, eachSource := range Fresp.SourceBk {
		fmt.Println("File ", eachSource.File)
		fmt.Println("Label ", eachSource.Label)
		fmt.Println("Type ", eachSource.Type)
	}
	return []Link{}

}
