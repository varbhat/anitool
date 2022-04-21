package main

import (
	"fmt"
	"os"
	"os/exec"
)

const USERAGENT string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36"

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
		referer = fmt.Sprintf(`--referrer="%s"`, l.Referer)
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf(`mpv %s %s %s`, fmt.Sprintf(`--user-agent="%s"`, USERAGENT), referer, fmt.Sprintf(`"%s"`, l.File)))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()
}
