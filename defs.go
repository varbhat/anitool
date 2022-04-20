package main

import (
	"fmt"
	"os"
	"os/exec"
)

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
	cmd := exec.Command("mpv", `--user-agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 12_3_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36"`, referer, l.File)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()
}
