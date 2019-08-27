// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Gotest process JSON stream generated by `go test -json` or `test2json`.
//
// Usage:
//
//	gotest [go test -json]
//
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/JulienVdG/tastevin/pkg/browser"
	"github.com/JulienVdG/tastevin/pkg/gotest"
	"github.com/JulienVdG/tastevin/pkg/gotestweb"
	"github.com/JulienVdG/tastevin/pkg/json2test"
	"github.com/JulienVdG/tastevin/pkg/testsuite"
	"github.com/JulienVdG/tastevin/pkg/xio"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: gotest [flags] <run|live|serve|gen> [go test -json]\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  run\trun the test and save outputs\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  live\trun the test, save outputs and serve them in browser\n")
	fmt.Fprintf(flag.CommandLine.Output(), "  serve\tserve the given test output in browser\n\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	flagL = flag.String("l", "logs/", "set the log directory, also passed to tastevin config")
	flagJ = flag.String("j", "results.json", "set JSON `filename` inside log directory")
	flagV = flag.Bool("v", false, "verbose test output (like go test -v)")
	flagS = flag.Bool("s", false, "silent (ie no test output)")
	flagP = flag.String("p", "", "dev: proxy to gotestweb (typical http://localhost:3000/ )")
	flagT = flag.String("t", "", "dev: test path to gotestweb (typical pkg/gotestweb/webapp/dist/ )")
)

func main() {
	var doRun, live, doServe bool
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}
	args := flag.Args()
	switch args[0] {
	case "run":
		doRun = true
	case "live":
		doRun = true
		live = true
		doServe = true
	case "serve":
		doServe = true
	default:
		usage()
	}

	var l io.WriteCloser
	serveDone := make(chan struct{}, 1)
	if doServe {
		l = serve(live, serveDone)
	}

	var err error
	if doRun {
		err = run(l, args[1:])
	}
	if live {
		fmt.Printf("gotest live: test done.\n")
	}

	if doServe {
		// Wait for http.ListenAndServe end
		<-serveDone
	}

	if err != nil {
		if _, ok := err.(gotest.RunError); !ok {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

func setHTTPHandlers() (string, error) {
	if *flagP != "" {
		rpURL, err := url.Parse(*flagP)
		if err != nil {
			return "", err
		}
		http.Handle("/", httputil.NewSingleHostReverseProxy(rpURL))
	} else if *flagT != "" {
		http.Handle("/", http.FileServer(http.Dir(*flagT)))
	} else {
		err := gotestweb.Handle()
		if err != nil {
			return "", err
		}
	}
	slug := filepath.Base(*flagL)
	prefix := "/" + slug + "/"

	http.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(*flagL))))
	return slug, nil
}

func serve(live bool, done chan struct{}) io.WriteCloser {
	var l io.WriteCloser
	slug, err := setHTTPHandlers()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if live {
		l = gotestweb.HandleLive()
	}
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
		done <- struct{}{}
	}()
	url := "http://localhost:8080/#"
	if live {
		url += slug + "/?live&summary=1"
	} else {
		url += slug + "/" + *flagJ + "?asciicast=" + slug + "&summary=0"
	}
	browser.Open(url)
	return l
}

func run(l io.WriteCloser, args []string) error {
	var c io.WriteCloser

	// Update env
	absdir, err := filepath.Abs(*flagL)
	if err != nil {
		return fmt.Errorf("error getting absolute path '%s': %v", *flagL, err)
	}
	err = testsuite.SetConfLogDir(absdir + "/")
	if err != nil {
		return fmt.Errorf("error updating config: %v", err)
	}

	if !*flagS {
		var h json2test.TestEventHandler
		if *flagV {
			h = json2test.NewVerboseHandler(os.Stdout)
		} else {
			h = json2test.NewSummaryHandler(os.Stdout)
		}

		c = json2test.NewConverter(h)
	}
	if *flagJ != "" {
		jsonfilename := filepath.Join(absdir, *flagJ)
		dir := filepath.Dir(jsonfilename)
		err := os.MkdirAll(dir, 0775)
		if err != nil {
			return fmt.Errorf("error creating directory '%s': %v", dir, err)
		}

		f, err := os.Create(jsonfilename)
		if err != nil {
			return fmt.Errorf("error creating file '%s': %v", *flagJ, err)
		}
		if l == nil {
			if *flagS {
				c = f
			} else {
				c = xio.MultiWriteCloser(c, f)
			}
		} else {
			if *flagS {
				c = xio.MultiWriteCloser(l, f)
			} else {
				c = xio.MultiWriteCloser(l, c, f)
			}
		}
	} else if *flagS {
		if l == nil {
			return errors.New("error -j is required in silent mode")
		}
		c = l
	}

	if len(args) == 0 {
		io.Copy(c, os.Stdin)
		c.Close()
		return nil
	}
	return gotest.Run(c, args)
}
