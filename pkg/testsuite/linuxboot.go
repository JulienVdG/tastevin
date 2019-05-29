// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testsuite

import (
	"testing"

	exp "github.com/google/goexpect"
)

// Linuxboot2urootBatcher follows the boot sequence of u-root to the shell prompt
var Linuxboot2urootBatcher []exp.Batcher = []exp.Batcher{
	&BExpTLog{
		L: "Matched u-root banner",
		R: "Welcome to u-root!",
		T: 40, // TODO make this time a parameter
	}, &BExpTLog{
		L: "Matched u-root prompt",
		R: "~/> ",
		T: 5,
	}}

// Linuxboot2uroot test the boot sequence of u-root to the shell prompt
func Linuxboot2uroot(t *testing.T, e *exp.GExpect) error {
	batcher := Linuxboot2urootBatcher
	res, err := e.ExpectBatch(batcher, 0)
	if err != nil {
		t.Errorf("Linuxboot2uroot: e.ExpectBatch(%v,_), err: %v, res: %q", batcher, err, res)

	}
	return nil
}
