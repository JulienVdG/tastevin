// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testsuite

import (
	"fmt"
	"regexp"
	"time"

	exp "github.com/google/goexpect"
	"google.golang.org/grpc/codes"
)

// TagLog adds logging to goexpect.OK
func TagLog(msg string) func() (exp.Tag, *exp.Status) {
	return func() (exp.Tag, *exp.Status) {
		// Don't use t.Log here, we want the event to be synchronous
		// to be useful as seek points for asciinema in gotest-web
		// however go test loose the test origin (as -json is a filter
		// not the internal, it has no way to scope the output in the
		// right test.)
		fmt.Printf("%s\n", msg)
		return exp.OKTag, exp.NewStatus(codes.OK, "state reached")
	}
}

// BExpTLog implements the Batcher interface for Expect commands with timeout and log on success
type BExpTLog struct {
	// R contains the Expect command regular expression.
	R string
	// T holds the Expect command timeout in seconds.
	T int
	// L holds the string to log once matched.
	L string
}

// Timeout returns the timeout in seconds.
func (betl *BExpTLog) Timeout() time.Duration {
	return time.Duration(betl.T) * time.Second
}

// Cmd returns the SwitchCase command(BatchSwitchCase).
func (betl *BExpTLog) Cmd() int {
	return exp.BatchSwitchCase
}

// Arg returns an empty string , not used for SwitchCase.
func (betl *BExpTLog) Arg() string {
	return ""
}

// Cases returns the Caser structure.
func (betl *BExpTLog) Cases() []exp.Caser {
	return []exp.Caser{&exp.Case{
		R: regexp.MustCompile(betl.R),
		T: TagLog(betl.L),
	}}
}
