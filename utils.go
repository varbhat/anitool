package main

import (
	"fmt"
	"net/http"
	"os/exec"
)

func finalURL(url string) (ret string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	return resp.Request.URL.String(), nil
}

func runMPV(exe string, url string, referer string) {
	if exe == "" {
		exe = "mpv"
	}
	if url == "" {
		return
	}
	if referer != "" {
		referer = fmt.Sprintf(`--referrer="%s"`, referer)
	}
	cmd := exec.Command(exe, url, referer)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}

func runVLC(exe string, url string, referer string) {
	if exe == "" {
		exe = "vlc"
	}
	if url == "" {
		return
	}
	if referer != "" {
		referer = fmt.Sprintf(`--http-referrer="%s"`, referer)
	}
	cmd := exec.Command(exe, url, "--play-and-exit", referer)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}
