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
		referer = fmt.Sprintf(`--referrer="%s"`, referer)
	}
	cmd := exec.Command("mpv", l.File, referer)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()
}
