// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gotest process JSON stream generated by `go test -json` or `test2json`.
//
// Usage:
//
//	gotestlive [go test -json]
//
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"

	"github.com/JulienVdG/tastevin/pkg/browser"
	"github.com/JulienVdG/tastevin/pkg/gotestweb"
)

func serve() {
	// TODO param for proxy
	if false {
		rpURL, err := url.Parse("http://localhost:3000/")
		if err != nil {
			log.Fatal(err)
		}
		http.Handle("/", httputil.NewSingleHostReverseProxy(rpURL))
	} else {
		http.Handle("/", http.FileServer(http.Dir("../gotest-web/gotest-web2/dist/")))
	}
	http.Handle("/single/", http.StripPrefix("/single/", http.FileServer(http.Dir("logs/"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func run(l io.WriteCloser, args []string) error {
	defer l.Close()
	cmd := exec.Command(args[0], args[1:]...)
	w := &countWriter{0, l}
	cmd.Stdout = w
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		if w.n > 0 {
			// Assume command printed why it failed.
		} else {
			fmt.Fprintf(l, "Error unable to run tests: %v\n", err)
			fmt.Printf("gotestlive: %v\n", err)
		}
		return err
	}
	return nil
}

type countWriter struct {
	n int64
	w io.Writer
}

func (w *countWriter) Write(b []byte) (int, error) {
	w.n += int64(len(b))
	return w.w.Write(b)
}

func main() {
	flag.Parse()
	l := gotestweb.HandleLive()
	errc := make(chan error, 1)
	go func() {
		serve()
		errc <- nil
	}()
	browser.Open("http://localhost:8080/#build?live&asciicast=single&summary")
	args := flag.Args()
	err := run(l, args)
	fmt.Printf("gotestlive: test done.\n")

	<-errc // Wait for serve end
	if err != nil {
		os.Exit(1)
	}
}
