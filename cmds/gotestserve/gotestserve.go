// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gotestserve start an http server to show go test -json in a browser
//
// Usage:
//
//	gotestserve [path_to/results.json]
//
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"

	"github.com/JulienVdG/tastevin/pkg/browser"
	"github.com/JulienVdG/tastevin/pkg/gotestweb"
)

var (
	flagP = flag.Bool("p", false, "proxy for dev")
	flagL = flag.Bool("l", false, "local files for prereleases")
)

func setHTTPHandlers(dir string) {
	// TODO param for proxy
	if *flagP {
		rpURL, err := url.Parse("http://localhost:3000/")
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/", httputil.NewSingleHostReverseProxy(rpURL))
	} else if *flagL {
		http.Handle("/", http.FileServer(http.Dir("../gotest-web/gotest-web2/dist/")))
	} else {
		err := gotestweb.Handle()
		if err != nil {
			log.Fatal(err)
		}
	}
	http.Handle("/single/", http.StripPrefix("/single/", http.FileServer(http.Dir(dir))))

}

func main() {
	flag.Parse()
	json := "logs/results.json"
	if flag.NArg() != 0 {
		json = flag.Arg(0)
	}
	dir := filepath.Dir(json)
	setHTTPHandlers(dir)
	errc := make(chan error, 1)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
		errc <- nil
	}()
	url := fmt.Sprintf("http://localhost:8080/#build?file=single/%s&asciicast=single&summary=0", filepath.Base(json))
	browser.Open(url)
	fmt.Printf("gotestserve: serving on %s.\n", url)

	<-errc // Wait for http.ListenAndServe end
}
