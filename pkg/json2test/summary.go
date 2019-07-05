// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json2test

import (
	"io"
	"strings"
)

type summaryHandler struct {
	w io.Writer // stdout for testing output stream
}

// NewSummaryHandler return a TestEventHandler that writes 'output' events to w.
func NewSummaryHandler(w io.Writer) TestEventHandler {
	return summaryHandler{w: w}
}

func (v summaryHandler) Handle(e TestEvent) {
	// drop test event (keep only package events)
	if e.Test != "" {
		return
	}

	if e.Action == "output" {
		// drop package PASS/FAIL summary
		if e.Output == "PASS\n" || e.Output == "FAIL\n" {
			return
		}
		// drop coverage output
		if strings.HasPrefix(e.Output, "coverage:") && strings.HasSuffix(e.Output, "% of statements\n") {
			return
		}
		_, err := v.w.Write([]byte(e.Output))
		if err != nil {
			LogWarn("json2test summary output: %s", err)
		}
	}
}
