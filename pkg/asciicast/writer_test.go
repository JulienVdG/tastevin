// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asciicast_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/asciicast"
)

type CallCounter struct {
	Name      string
	CallCount int
}

type ExpectWriter struct {
	CallCounter
	Content []byte
	Result  error
}

type ExpectCloser struct {
	CallCounter
}

type ExpectWriteCloser struct {
	ExpectWriter
	ExpectCloser
}

func (ew *ExpectWriter) Write(p []byte) (n int, err error) {
	ew.CallCount++
	cmp := bytes.Compare(ew.Content, p)
	if cmp != 0 {
		ew.Result = fmt.Errorf("mismatched %s write, wanted '%v', got '%v'", ew.Name, string(ew.Content), string(p))
	}
	return len(p), nil
}

func (ec *ExpectCloser) Close() error {
	ec.CallCount++
	return nil
}

func (cc *CallCounter) NotCalled() (err error) {
	if cc.CallCount != 0 {
		err = fmt.Errorf("%s called while not expected, got %v", cc.Name, cc.CallCount)
	}
	cc.CallCount = 0
	return
}

func (cc *CallCounter) CalledOnce() (err error) {
	if cc.CallCount == 0 {
		err = fmt.Errorf("%s not called, expected once", cc.Name)
	} else if cc.CallCount != 1 {
		err = fmt.Errorf("%s called more than once, got %v", cc.Name, cc.CallCount)
	}
	cc.CallCount = 0
	return
}

func Assert(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

var testtime time.Time

func TestWriter(t *testing.T) {
	initialtime := time.Unix(0, 0)
	// Mock the time getter
	testtime = initialtime
	asciicast.Now = func() time.Time { return testtime }

	var script ExpectWriteCloser
	script.ExpectWriter.Name = "script writer"
	script.ExpectCloser.Name = "script closer"

	script.Content = []byte(fmt.Sprintf("{\"version\": 2, \"width\": 80, \"height\": 24, \"timestamp\": %d}\n", testtime.Unix()))

	sr := asciicast.NewWriter(&script)

	Assert(t, script.ExpectWriter.Result)
	Assert(t, script.ExpectWriter.CalledOnce())
	Assert(t, script.ExpectCloser.NotCalled())

	var tests = []struct {
		content string
		delay   float64
	}{
		{"Send some char\n", 1.05},
		{"more lines\n...\neven more", 0.01},
		{"minimal timing", 0.000001},
	}
	for _, test := range tests {
		testtime = testtime.Add(time.Duration(test.delay * float64(time.Second)))
		to := testtime.Sub(initialtime)
		script.Content = []byte(fmt.Sprintf("[%.06f, \"o\", %s]\n", to.Seconds(), strconv.QuoteToASCII(test.content)))

		sr.Write([]byte(test.content))

		Assert(t, script.ExpectWriter.Result)
		Assert(t, script.ExpectWriter.CalledOnce())
		Assert(t, script.ExpectCloser.NotCalled())
	}

	script.Content = []byte(fmt.Sprintf("\nScript done on %s [<end>]\n", testtime.Format(time.RFC3339)))

	sr.Close()

	Assert(t, script.ExpectWriter.NotCalled())
	Assert(t, script.ExpectCloser.CalledOnce())
}
