// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testsuite

import (
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
		timeout: 20 * time.Second,
	}, {
		name:    "Match prompt",
		re:      regexp.MustCompile("~/> "),
		timeout: 1 * time.Second,
	}}
	for _, tst := range tests {
		out, _, err := e.Expect(tst.re, tst.timeout)
		if err != nil {
			t.Errorf("%s: Expect(%q,%v), err: %v, out: %q", tst.name, tst.re.String(), tst.timeout, err, out)
			continue
		}
		//t.Log(out)
	}

	return nil
}
