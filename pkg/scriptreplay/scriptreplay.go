// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package scriptreplay handle logs compatible with the scriptreplay command
package scriptreplay

import (
	"fmt"
	"io"
	"os"
	"time"
)

type sr struct {
	script  io.WriteCloser
	timing  io.WriteCloser
	oldtime time.Time
}

// Now function provide time default to time.Now
var Now = time.Now

// NewWriter provide a WriteCloser that writes script format to script, timing
func NewWriter(script, timing io.WriteCloser) io.WriteCloser {
	t := Now()
	sr := sr{script: script, timing: timing, oldtime: t}

	msg := fmt.Sprintf("Script started on %s [<not executed on terminal>]\n", t.Format(time.RFC3339))
	sr.script.Write([]byte(msg))

	return &sr
}

// NewFileWriter provide a WriteCloser that writes script format to file
func NewFileWriter(scriptName, timingName string) (io.WriteCloser, error) {
	fo, err := os.Create(scriptName)
	if err != nil {
		return nil, fmt.Errorf("Error creating file '%s': %v", scriptName, err)
	}

	ft, err := os.Create(timingName)
	if err != nil {
		fo.Close()
		return nil, fmt.Errorf("Error creating file '%s': %v", timingName, err)
	}

	return NewWriter(fo, ft), nil

}

// Write implement io.Writer interface
func (sr *sr) Write(p []byte) (n int, err error) {
	t := Now()
	diff := t.Sub(sr.oldtime)
	sr.oldtime = t
	n, err = sr.script.Write(p)
	msg := fmt.Sprintf("%.06f %d\n", diff.Seconds(), n)
	sr.timing.Write([]byte(msg))
	return
}

// Close implement io.Closer interface
func (sr *sr) Close() error {
	t := Now()
	msg := fmt.Sprintf("\nScript done on %s [<end>]\n", t.Format(time.RFC3339))
	sr.script.Write([]byte(msg))

	var err error
	errt := sr.timing.Close()
	errl := sr.script.Close()
	if errt != nil && errl != nil {
		err = fmt.Errorf("error closing both timing and script logger %v, %v", errt, errl)
	} else if errt != nil {
		err = fmt.Errorf("error closing timing logger %v", errt)
	} else if errl != nil {
		err = fmt.Errorf("error closing script logger %v", errl)
	}
	return err
}
