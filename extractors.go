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
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func finalURL(url string) (ret string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	return resp.Request.URL.String(), nil
}

func gogoStreamLinks(aid string, epno string) (Ret chan Link) {
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

func (l *Link) Play() {
	if l.File == "" {
		return
	}
	referer := ""
	if l.Referer != "" {
		referer = fmt.Sprintf(`--referrer="%s"`, referer)
	}
	cmd := exec.Command("mpv", l.File, referer)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func GoGoCDN(iurl string) (Ret []Link) {
	type Fresponse struct {
		Source   []Link `json:"source"`
		SourceBk []Link `json:"source_bk"`
	}

	secret_key := "93106165734640459728346589106791"
	second_key := "97952160493714852094564712118349"
	iv := "8244002440089157"

	if strings.Contains(iurl, "streaming.php") {
		iUrl, err := url.Parse(iurl)
		if err != nil {
			return []Link{}
		}

		response, err := http.Get(iurl)
		if err != nil {
			return []Link{}
		}

		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			return []Link{}
		}

		encrypted := doc.Find("script[data-name='episode']").AttrOr("data-value", "")

		decrypted, err := aes256decrypt(encrypted, []byte(secret_key), []byte(iv), aes.BlockSize)
		if err != nil {
			fmt.Println("err ", err)
		}
		vididno := strings.IndexRune(decrypted, '&')
		vidid := decrypted[:vididno]
		vidend := decrypted[vididno:]

		encryptedvidid, err := aes256encrypt([]byte(vidid), []byte(secret_key), []byte(iv), aes.BlockSize)
		if err != nil {
			fmt.Println("err ", err)
		}

		ajax_url := fmt.Sprintf("https://%s/encrypt-ajax.php?id=%s%s&alias=%s", iUrl.Host, encryptedvidid, vidend, vidid)

		req, err := http.NewRequest("POST", ajax_url, bytes.NewBuffer([]byte(ajax_url)))
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Set("Referer", iUrl.Host)
		req.Header.Set("x-requested-with", "XMLHttpRequest")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		var ExtractData struct {
			Data string `json:"data"`
		}

		err = json.Unmarshal(body, &ExtractData)
		if err != nil {
			fmt.Println("Err ", err)
		}
		data, err := aes256decrypt(ExtractData.Data, []byte(second_key), []byte(iv), aes.BlockSize)
		if err != nil {
			fmt.Println("Err ", err)
			return
		}

		var Fresp Fresponse

		if err := json.Unmarshal([]byte(data), &Fresp); err != nil {
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

		for each := range Ret {
			fmt.Printf("%v", each)
		}
		return
	}
	return
}

func Fplayer(iurl string) (Ret []Link) {
	type FplayerResp struct {
		Data []Link `json:"data"`
	}
	req, err := http.NewRequest("POST", strings.Replace(iurl, "/v/", "/api/source/", -1), nil)
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

func StreamSB(iurl string) (Ret []Link) {
	var SSBresp struct {
		SSD struct {
			File   string `json:"file"`
			Title  string `json:"title"`
			Backup string `json:"backup"`
		} `json:"stream_data"`
	}
	ssburl := "https://sbplay2.com"
	jsonlink := ssburl + "/sources40/7361696b6f757c7c%s7c7c7361696b6f757c7c73747265616d7362/7361696b6f757c7c363136653639366436343663363136653639366436343663376337633631366536393664363436633631366536393664363436633763376336313665363936643634366336313665363936643634366337633763373337343732363536313664373336327c7c7361696b6f757c7c73747265616d7362"

	u, err := url.Parse(iurl)
	if err != nil {
		return
	}
	params := strings.Split(u.Path, "/")
	paramlen := len(params)
	if paramlen < 2 || params[1] != "e" {
		return

	}
	req, err := http.NewRequest("GET", fmt.Sprintf(jsonlink, hex.EncodeToString([]byte(params[2]))), nil)
	if err != nil {
		return
	}
	req.Header.Set("watchsb", "streamsb")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer response.Body.Close()
	respb, err := io.ReadAll(response.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(respb, &SSBresp)
	if err != nil {
		return
	}
	Ret = append(Ret, Link{
		File:    SSBresp.SSD.File,
		Label:   SSBresp.SSD.Title,
		Type:    "hls",
		Referer: iurl,
	})
	Ret = append(Ret, Link{
		File:    SSBresp.SSD.Backup,
		Label:   SSBresp.SSD.Title + " (Backup)",
		Type:    "hls",
		Referer: iurl,
	})
	return
}
