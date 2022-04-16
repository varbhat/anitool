package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func gogoStreamLinks(gogobaseurl string, aid string, epno string) (ret chan string) {
	ret = make(chan string)
	go func() {
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
	}()
	return
}

func aes256encrypt(plaintext []byte, key []byte, iv []byte, blockSize int) (ret string, err error) {
	defer func() {
		if r := recover(); r != nil {
			ret = ""
			err = errors.New("aes256 encryption failed")
			return
		}
	}()
	bPlaintext := pKCS5Padding(plaintext, blockSize, len(plaintext))
	block, err := aes.NewCipher(key)
	if err != nil {
		ret = ""
		return
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)
	ret = base64.StdEncoding.EncodeToString(ciphertext)
	return
}

func aes256decrypt(plaintext string, key []byte, iv []byte, blocksize int) (ret string, err error) {
	defer func() {
		if r := recover(); r != nil {
			ret = ""
			err = errors.New("aes256 decryption failed")
			return
		}
	}()
	bPlaintext, err := base64.StdEncoding.DecodeString(plaintext)
	if err != nil {
		ret = ""
		err = errors.New("aes256 decryption failed")
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		ret = ""
		err = errors.New("aes256 decryption failed")
		return
	}
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)
	ciphertext = pkCS7unpad(ciphertext, blocksize)
	ret = string(ciphertext)
	return
}

func pKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pkCS7unpad(padded []byte, size int) []byte {
	if len(padded)%size != 0 {
		return nil
	}

	bufLen := len(padded) - int(padded[len(padded)-1])
	buf := make([]byte, bufLen)
	copy(buf, padded[:bufLen])
	return buf
}

type Link struct {
	File    string
	Label   string
	Type    string
	Referer string
}

func GoGoCDN(iurl string) (Ret []Link) {
	type Fresponse struct {
		Source   []Link `json:"source"`
		SourceBk []Link `json:"source_bk"`
	}
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
	response, err := http.Get(iurl)
	if err != nil {
		return []Link{}
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return []Link{}
	}

	encrypted := doc.Find("script[data-name='crypto']").AttrOr("data-value", "")
	iv := doc.Find("script[data-name='ts']").AttrOr("data-value", "")
	if encrypted == "" || iv == "" {
		secret_key = "3235373436353338353932393338333936373634363632383739383333323838"
		iv = "34323036393133333738303038313335"
		rtime = "69420691337800813569"
		key, err := hex.DecodeString(secret_key)
		if err != nil {
			return
		}
		iv, err := hex.DecodeString(iv)
		if err != nil {
			return
		}
		ajax, err = aes256encrypt([]byte(video_id), key, iv, aes.BlockSize)
		if err != nil {
			return
		}
	} else {
		rtime = "00000000000000000000"
		secret_key, err = aes256decrypt(encrypted, []byte(iv+iv), []byte(iv), aes.BlockSize)
		if err != nil {
			return
		}
		ajax, err = aes256encrypt([]byte(video_id), []byte(secret_key), []byte("0000000000000000"), aes.BlockSize)
		if err != nil {
			return
		}
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
		var er Link
		er.File = eachSource.File
		er.Label = eachSource.Label
		er.Type = eachSource.Type
		er.Referer = iurl
		Ret = append(Ret, er)
	}

	for _, eachSource := range Fresp.SourceBk {
		var er Link
		er.File = eachSource.File
		er.Label = eachSource.Label + "(Backup)"
		er.Type = eachSource.Type
		er.Referer = iurl
		Ret = append(Ret, er)
	}
	return

}

func Fplayer(iurl string) (Ret []Link) {
	type FplayerResp struct {
		Data []Link `json:"data"`
	}
	apiurl := strings.Replace(iurl, "/v/", "/api/source/", -1)
	req, err := http.NewRequest("POST", apiurl, nil)
	if err != nil {
		return
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()
	var Fresp FplayerResp
	respb, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	if err := json.Unmarshal(respb, &Fresp); err != nil {
		return
	}

	for _, eachSource := range Fresp.Data {
		eachSource.Referer = iurl
	}
	Ret = Fresp.Data
	return
}

func StreamSB(iurl string) string {
	ssburl := "https://sbplay2.com"
	jsonlink := ssburl + "/sources40/7361696b6f757c7c%s7c7c7361696b6f757c7c73747265616d7362/7361696b6f757c7c363136653639366436343663363136653639366436343663376337633631366536393664363436633631366536393664363436633763376336313665363936643634366336313665363936643634366337633763373337343732363536313664373336327c7c7361696b6f757c7c73747265616d7362"

	u, err := url.Parse(iurl)
	if err != nil {
		return ""
	}
	params := strings.Split(u.Path, "/")
	paramlen := len(params)
	if paramlen < 2 {
		return ""

	}
	if params[1] != "e" {
		return ""
	}
	sourceid := hex.EncodeToString([]byte(params[2]))
	jsonApi := fmt.Sprintf(jsonlink, sourceid)
	req, err := http.NewRequest("GET", jsonApi, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("watchsb", "streamsb")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)

	}
	defer response.Body.Close()
	respb, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	fmt.Println(respb)
	return ""
}
