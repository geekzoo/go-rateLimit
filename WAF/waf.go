package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	goratelimit "github.com/geekzoo/go-rateLimit"
)

var (
	proxy    = &httputil.ReverseProxy{Director: director}
	director = func(req *http.Request) {
		req.Header.Add("X-Forwarded-Host", req.Host)
		req.Header.Add("X-Origin-Host", origin.Host)
		req.URL.Scheme = "http"
		req.URL.Host = origin.Host
	}
	origin, _ = url.Parse("http://192.168.49.2/")
	listen    = "0.0.0.0:3128"
	rate      = 50
)

func main() {

	if len(os.Getenv("LISTEN")) >= 0 {
		log.Println("LISTEN empty using: ", listen)
	} else {
		listen = os.Getenv("LISTEN")
	}
	if len(os.Getenv("ORIGIN")) >= 0 {
		log.Println("ORIGIN empty using: ", origin)
	} else {
		origin, _ = url.Parse(os.Getenv("ORIGIN"))
	}
	os.Getenv("ORIGIN")
	os.Getenv("LISTEN")
	os.Getenv("RATE")
	os.Getenv("RATE_TIMER")

	http.HandleFunc("/", extern)
	http.HandleFunc("/status", externPath)
	log.Fatal(http.ListenAndServe(listen, nil))
}

func extern(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Via", "v8-rate")
	token := r.Header.Get("authorization")
	if len(token) == 0 {
		proxy.ServeHTTP(w, r)
	} else {
		allowed := goratelimit.RateLimit(token, rate)
		if !allowed {
			w.Header().Add("X-Rate-Limit", fmt.Sprintf("%v", rate))
			w.WriteHeader(429)
		} else {
			proxy.ServeHTTP(w, r)
		}
	}
}

func externPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Via", "v8-rate")
	w.Header().Add("X-Extern-Path", "NEW")
	rate = 10
	origin, _ = url.Parse("http://localhost:8000/status")
	fmt.Println(origin)
	token := r.Header.Get("token")
	if len(token) == 0 {
		proxy.ServeHTTP(w, r)
	} else {
		allowed := goratelimit.RateLimit(token, rate)
		if !allowed {
			w.Header().Add("X-Rate-Limit", fmt.Sprintf("%v", rate))
			w.WriteHeader(429)
		} else {
			proxy.ServeHTTP(w, r)
		}
	}
}
