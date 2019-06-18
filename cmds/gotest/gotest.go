// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on test2json/main.go
// Copyright 2017 The Go Authors. All rights reserved.

// Gotest process JSON stream generated by `go test -json` or `test2json`.
//
// Usage:
//
//	gotest [go test -json]
//
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/JulienVdG/tastevin/pkg/json2test"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gotest [-j file] [go test -json]\n")
	os.Exit(2)
}

var (
	flagJ = flag.String("j", "", "save JSON to `file`")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	c := json2test.NewConverter(json2test.NewVerboseHandler(os.Stdout))
	if *flagJ != "" {
		f, err := os.Create(*flagJ)
		if err != nil {
			fmt.Printf("Error creating file '%s': %v", *flagJ, err)
			os.Exit(1)
		}
		defer f.Close()
		c = io.MultiWriter(c, f)
	}

	if flag.NArg() == 0 {
		io.Copy(c, os.Stdin)
	} else {
		args := flag.Args()
		cmd := exec.Command(args[0], args[1:]...)
		w := &countWriter{0, c}
		cmd.Stdout = w
		cmd.Stderr = w
		if err := cmd.Run(); err != nil {
			if w.n > 0 {
				// Assume command printed why it failed.
			} else {
				fmt.Printf("test2json: %v\n", err)
			}
			os.Exit(1)
		}
	}
}

type countWriter struct {
	n int64
	w io.Writer
}

func (w *countWriter) Write(b []byte) (int, error) {
	w.n += int64(len(b))
	return w.w.Write(b)
}
