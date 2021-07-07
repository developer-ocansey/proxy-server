package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
)

type proxy struct {
	rmt *url.URL
}

func (p *proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	client := &http.Client{}

	req, err := http.NewRequest(req.Method, p.rmt.String(), req.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(wr, "Server Error", http.StatusInternalServerError)
		log.Fatal("ServeHTTP:", err)
	}
	defer resp.Body.Close()

	log.Println(p.rmt, " ", resp.Status)

	wr.WriteHeader(resp.StatusCode)
	io.Copy(wr, resp.Body)
}

func main() {
	var addr = flag.String("addr", ":9090", "The addr of the application.")
	var rmt = flag.String("rmt", "localhost:8080/test", "Remote address.")
	flag.Parse()

	url, err := url.Parse(*rmt)
	if err != nil {
		log.Fatal(err)
	}
	handler := &proxy{
		rmt: url,
	}

	log.Println("Starting proxy server on", *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
