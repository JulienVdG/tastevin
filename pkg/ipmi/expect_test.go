// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/ipmi"
	"github.com/JulienVdG/tastevin/pkg/scriptreplay"
	"github.com/JulienVdG/tastevin/pkg/testsuite"
	exp "github.com/google/goexpect"
)

// Run with env TASTEVIN_IPMI='{"Hostname":"10.0.3.208","Username":"USERID","Password":"PASSW0RD","Interface":"lanplus","Path":"ipmitool"}'
func TestLinuxboot2uroot(t *testing.T) {
	c, err := ipmi.ConfigFromEnv()
	if err != nil {
		t.Skipf("IPMI test is skipped unless TASTEVIN_IPMI is set (%v)", err)
	}
	r, err := ipmi.NewRemote(c)
	if err != nil {
		t.Fatal(err)
	}

	logdir := filepath.Join("testdata", "log")
	err = os.MkdirAll(logdir, 0775)
	if err != nil {
		t.Fatalf("TeeReplay failed: %v", err)
	}

	sr, err := scriptreplay.NewFileWriter(filepath.Join(logdir, "IPMILinuxboot2uroot.log"), filepath.Join(logdir, "IPMILinuxboot2uroot.tim"))
	if err != nil {
		t.Fatalf("TeeReplay failed: %v", err)
	}

	err = r.Open()
	if err != nil {
		t.Error(err)
	}

	err = r.PowerUp()
	if err != nil {
		t.Error(err)
	}

	// Don't connect SOL immediately after power change
	time.Sleep(5 * time.Second)

	e, _, err := r.Spawn(1*time.Second, exp.PartialMatch(true), exp.Tee(sr) /* exp.DebugCheck(nil), exp.Verbose(true)*/)
	if err != nil {
		t.Fatalf("Spawn failed: %v", err)
	}

	err = testsuite.Linuxboot2uroot(t, e)
	if err != nil {
		t.Fatalf("Linuxboot2uroot returned: %v", err)
	}

	err = r.PowerDown()
	if err != nil {
		t.Error(err)
	}

	err = r.Close()
	if err != nil {
		t.Error(err)
	}

	err = sr.Close()
	if err != nil {
		t.Errorf("sr close: %v", err)
	}
}
