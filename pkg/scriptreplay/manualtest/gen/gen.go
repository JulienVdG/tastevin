// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/JulienVdG/tastevin/pkg/scriptreplay"
)

var testtime time.Time

func main() {
	// Mock the time getter
	testtime = time.Unix(0, 0)
	scriptreplay.Now = func() time.Time { return testtime }
	fmt.Println("Generating `output.txt` and `output.tim`.")

	// open output file
	fo, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}

	// open output file
	ft, err := os.Create("output.tim")
	if err != nil {
		panic(err)
	}

	sr := scriptreplay.NewWriter(fo, ft)

	var tests = []struct {
		content string
		delay   float64
	}{
		{"Send some char\n", 1.05},
		{"more lines\n...\neven more\n", 0.01},
		{"minimal timing\n", 0.000001},
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
		{"\n", 0.5},
		{"Done.\n", 0.000001},
	}
	for _, test := range tests {
		testtime = testtime.Add(time.Duration(test.delay * float64(time.Second)))

		sr.Write([]byte(test.content))

	}

	sr.Close()
	fmt.Println("Test with `scriptreplay output.txt -toutput.tim`.")
}
