// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"time"

	"github.com/JulienVdG/tastevin/pkg/asciicast"
)

var testtime time.Time

func main() {
	// Mock the time getter
	testtime = time.Unix(0, 0)
	asciicast.Now = func() time.Time { return testtime }
	fmt.Println("Generating `output.cast`.")

	ar, err := asciicast.NewFileWriter("output.cast")
	if err != nil {
		panic(err)
	}

	var tests = []struct {
		content string
		delay   float64
	}{
		{"Send some char\r\n", 1.05},
		{"more lines\r\n...\r\neven more\r\n", 0.01},
		{"minimal timing\r\n", 0.000001},
		{"O", 0.5},
		{"n", 0.5},
		{"e", 0.5},
		{" ", 0.5},
		{"b", 0.5},
		{"y", 0.5},
		{" ", 0.5},
		{"o", 0.5},
		{"n", 0.5},
		{"e", 0.5},
		{"!", 0.5},
		{"\r\n", 0.5},
		{"\xff\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f", 0.01},
		{"\x1b[31;1mRED\x1b[01;34mBlue\x1b[00mNormal\n", 0.001},
		{"Done.\r\n", 0.000001},
	}
	for _, test := range tests {
		testtime = testtime.Add(time.Duration(test.delay * float64(time.Second)))

		ar.Write([]byte(test.content))

	}

	ar.Close()
	fmt.Println("Test with `asciinema play output.cast`.")
}
