// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json2test

import (
	"io"
)

type verboseHandler struct {
	w io.Writer // stdout for testing output stream
}

// NewVerboseHandler return a TestEventHandler that writes 'output' events to w.
func NewVerboseHandler(w io.Writer) TestEventHandler {
	return verboseHandler{w: w}
}

func (v verboseHandler) Handle(e TestEvent) {
	if e.Action == "output" {
		_, err := v.w.Write([]byte(e.Output))
		if err != nil {
			LogWarn("json2test verbose output: %s", err)
		}
	}
}
