// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asciicast_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/asciicast"
	"github.com/JulienVdG/tastevin/pkg/xio/iotest"
)

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

	var script iotest.ExpectWriteCloser
	script.ExpectWriter.Name = "script writer"
	script.ExpectCloser.Name = "script closer"

	script.Content = []byte(fmt.Sprintf("{\"version\":2,\"width\":80,\"height\":24,\"timestamp\":%d}\n", testtime.Unix()))

	sr, err := asciicast.NewWriter(&script)
	Assert(t, err)

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
		// the QuoteToASCII is not JSON correct but OK for the above
		// strings, the encoding is tested manually on the web frontend
		// with asciinema-player
		script.Content = []byte(fmt.Sprintf("[%g,\"o\",%s]\n", to.Seconds(), strconv.QuoteToASCII(test.content)))

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
