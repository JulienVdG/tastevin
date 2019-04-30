// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testsuite_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
	exp "github.com/google/goexpect"
)

func TestLinuxboot2uroot(t *testing.T) {
	srv := []exp.Batcher{
		&exp.BSnd{`
2019/03/26 13:50:05 Welcome to u-root!
                              _
   _   _      _ __ ___   ___ | |_
  | | | |____| '__/ _ \ / _ \| __|
  | |_| |____| | | (_) | (_) | |_
   \__,_|    |_|  \___/ \___/ \__|

`},
		&exp.BSnd{`
~/> `},
	}

	logdir := filepath.Join("testdata", "log")
	err := os.MkdirAll(logdir, 0775)
	if err != nil {
		t.Fatalf("TeeReplay failed: %v", err)
	}

	opts, warn := testsuite.ExpectOptions("")
	if warn != nil {
		t.Log(warn)
	}
	// Force it to be verbose
	opts = append(opts, exp.DebugCheck(nil), exp.Verbose(true))

	e, ech, err := exp.SpawnFake(srv, 1*time.Second, opts...)
	if err != nil {
		t.Fatalf("SpawnFake failed: %v", err)
	}

	err = testsuite.Linuxboot2uroot(t, e)
	if err != nil {
		t.Errorf("Linuxboot2uroot returned: %v", err)
	}

	err = e.Close()
	if err != nil {
		t.Errorf("error closing SpawnFake: %v", err)
	}

	// make sure the expect session is done and screenlog are closed
	<-ech
}
