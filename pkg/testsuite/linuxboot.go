// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testsuite

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	exp "github.com/google/goexpect"
)

// Linuxboot2uroot test the boot sequence of u-root to the shell prompt
func Linuxboot2uroot(t *testing.T, e *exp.GExpect) error {

	tests := []struct {
		name    string
		fail    bool
		timeout time.Duration
		re      *regexp.Regexp
	}{{
		name:    "Match banner",
		re:      regexp.MustCompile("Welcome to u-root!"),
		timeout: 40 * time.Second, // TODO make this time a parameter
	}, {
		name:    "Match prompt",
		re:      regexp.MustCompile("~/> "),
		timeout: 5 * time.Second,
	}}
	for _, tst := range tests {
		out, _, err := e.Expect(tst.re, tst.timeout)
		if err != nil {
			t.Errorf("%s: Expect(%q,%v), err: %v, out: %q", tst.name, tst.re.String(), tst.timeout, err, out)
			continue
		}
		//t.Log(out)
		// Don't use t.Log here, we want the event to be synchronous
		// to be useful as seek points for asciinema in gotest-web
		// however go test loose the test origin (as -json is a filter
		// not the internal, it has no way to scope the output in the
		// right test.)
		fmt.Printf("%s done\n", tst.name)
	}

	return nil
}
