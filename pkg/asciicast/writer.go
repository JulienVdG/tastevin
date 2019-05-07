// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package asciicast handle logs compatible with the asciinema command
//
// asciicast file format (version 2) reference:
// https://github.com/asciinema/asciinema/blob/master/doc/asciicast-v2.md
package asciicast

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type ar struct {
	cast      io.WriteCloser
	starttime time.Time
	logclose  string
}

type header struct {
	Version   int   `json:"version"`
	Width     int   `json:"width"`
	Height    int   `json:"height"`
	Timestamp int64 `json:"timestamp"`
}

// Now function provide time default to time.Now
var Now = time.Now

// NewWriter provide a WriteCloser that writes asciicast format to cast
func NewWriter(cast io.WriteCloser) (io.WriteCloser, error) {
	t := Now()
	ar := ar{cast: cast, starttime: t}

	h := header{Version: 2, Width: 80, Height: 24, Timestamp: t.Unix()}
	msg, err := json.Marshal(&h)
	if err != nil {
		return nil, err
	}
	msg = append(msg, '\n')
	ar.cast.Write(msg)

	return &ar, nil
}

// NewFileWriter provide a WriteCloser that writes asciicast format to file
func NewFileWriter(castName string) (io.WriteCloser, error) {
	f, err := os.Create(castName)
	if err != nil {
		return nil, fmt.Errorf("Error creating file '%s': %v", castName, err)
	}
	wc, err := NewWriter(f)
	if err != nil {
		return nil, err
	}
	ar := wc.(*ar)
	fmt.Printf("*** Asciicast '%s' start\n", castName)
	ar.logclose = fmt.Sprintf("*** Asciicast '%s' end", castName)

	return wc, nil

}

// Write implement io.Writer interface
func (ar *ar) Write(p []byte) (n int, err error) {
	t := Now()
	to := t.Sub(ar.starttime)
	a := make([]interface{}, 3)
	a[0] = to.Seconds()
	a[1] = "o"
	a[2] = string(p)
	msg, err := json.Marshal(&a)
	if err != nil {
		return 0, err
	}
	msg = append(msg, '\n')

	_, err = ar.cast.Write([]byte(msg))
	return len(p), err
}

// Close implement io.Closer interface
func (ar *ar) Close() error {
	if ar.logclose != "" {
		fmt.Println(ar.logclose)
	}
	err := ar.cast.Close()
	if err != nil {
		return fmt.Errorf("error closing cast %v", err)
	}
	return nil
}
