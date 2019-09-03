// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Based on test2json/main.go
// Copyright 2017 The Go Authors. All rights reserved.

package main

import (
	"fmt"
	"io"
	"os/exec"
)

// Run execute the command in args and sends its standard and error output to the output writer
func Run(output io.WriteCloser, args []string) error {
	defer output.Close()
	cmd := exec.Command(args[0], args[1:]...)
	w := &countWriter{0, output}
	cmd.Stdout = w
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		if w.n > 0 {
			// Assume command printed why it failed.
		} else {
			fmt.Fprintf(output, "Error unable to run tests: %v\n", err)
			fmt.Printf("gotest: %v\n", err)
		}
		return RunError{err}
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

// RunError wraps error returned by Run to avoid printing it twice
type RunError struct {
	err error
}

func (r RunError) Error() string {
	return r.err.Error()
}
