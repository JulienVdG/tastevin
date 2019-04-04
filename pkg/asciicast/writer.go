// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package asciicast handle logs compatible with the asciinema command
package asciicast

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type ar struct {
	cast      io.WriteCloser
	starttime time.Time
}

// Now function provide time default to time.Now
var Now = time.Now

// NewWriter provide a WriteCloser that writes asciicast format to cast
func NewWriter(cast io.WriteCloser) io.WriteCloser {
	t := Now()
	ar := ar{cast: cast, starttime: t}

	// TODO use json
	msg := fmt.Sprintf("{\"version\": 2, \"width\": 80, \"height\": 24, \"timestamp\": %d}\n", t.Unix())
	ar.cast.Write([]byte(msg))

	return &ar
}

// NewFileWriter provide a WriteCloser that writes asciicast format to file
func NewFileWriter(castName string) (io.WriteCloser, error) {
	f, err := os.Create(castName)
	if err != nil {
		return nil, fmt.Errorf("Error creating file '%s': %v", castName, err)
	}

	return NewWriter(f), nil

}

// Write implement io.Writer interface
func (ar *ar) Write(p []byte) (n int, err error) {
	t := Now()
	to := t.Sub(ar.starttime)
	str := strconv.QuoteToASCII(string(p))
	msg := fmt.Sprintf("[%.06f, \"o\", %s]\n", to.Seconds(), str)
	_, err = ar.cast.Write([]byte(msg))
	return len(p), err
}

// Close implement io.Closer interface
func (ar *ar) Close() error {
	err := ar.cast.Close()
	if err != nil {
		return fmt.Errorf("error closing cast %v", err)
	}
	return nil
}
