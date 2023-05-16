package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var i = 0

func csum(in string) string {
	cs := 0
	for _, c := range in {
		if c == '\n' {
			continue
		}
		cs += int(c) + 1
	}

	return in + "\n" + strconv.Itoa(cs)
}

func handler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("not flushable")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	req := strings.Split(string(body), "\r\n")
	log.Println("request", req)

	if len(req) < 6 {
		return
	}

	resp := ""

	log.Println("req type", req[6])

	if req[6] == "U" {
		// text that shows in the ui + checksum
		resp = "No updates available."
	} else if req[6] == "PCTCIL" {
		// i have 0 clue what this means
		resp = "1\n5\n9"
	} else if req[6] == "C" {
		// 100 days left + checksum
		resp = "100"
	} else if req[6] == "ZVTAP2" {
		// this too
		resp = "1\n0 4096 1 0"
	} else if req[6] == "W2" {
		log.Println("not implemented")
		resp = "???"
	}

	resp = "START\n" + csum(resp) + "\nTRATS\n"
	resp = strings.ReplaceAll(resp, "\n", "\r\n")

	w.Header().Set("Cache-Control", "private")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Server", "Microsoft-IIS/10.0")
	w.Header().Set("Set-Cookie", "ASP.NET_SessionId=yaic0evvwudakk2sjwlmnyxi; path=/; HttpOnly; SameSite=Lax")
	w.Header().Set("X-ApsNet-Version", "4.0.30319")
	w.Header().Set("X-Powered-By", "ASP.NET")
	w.Write([]byte(resp))
	flusher.Flush()
}

func main() {
	panic(http.ListenAndServe(":80", http.HandlerFunc(handler)))
}
