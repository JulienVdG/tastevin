// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testsuite

import (
	"runtime"
	"strings"
)

func last(s string) string {
	l := strings.Split(s, ".")
	return l[len(l)-1]
}

func longFilename(s string) string {
	sanitize := func(r rune) rune {
		switch r {
		case '!', '$', '&', '\'', '(', ')', '*', '+', ',', ';', '=', ':', '[', ']', '<', '>', '"', '~', '?', '@':
			return -1
		case '.':
			return '-'
		case '/', '\\':
			return '_'

		}
		return r
	}
	return strings.Map(sanitize, s)
}

func callerName(depth int, longName bool) string {
	// Use the test name as the serial log's file name.
	pc, _, _, ok := runtime.Caller(depth)
	if !ok {
		panic("runtime caller failed")
	}
	f := runtime.FuncForPC(pc)
	if longName {
		return longFilename(f.Name())
	}
	return last(f.Name())
}
