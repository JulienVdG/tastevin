// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi_test

import (
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/ipmi"
	"github.com/JulienVdG/tastevin/pkg/testsuite"
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

	opts, warn := testsuite.ExpectOptions("")
	if warn != nil {
		t.Log(warn)
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

	e, _, err := r.Spawn(1*time.Second, opts...)
	if err != nil {
		t.Fatalf("Spawn failed: %v", err)
	}

	batcher := testsuite.Linuxboot2urootBatcher
	res, err := e.ExpectBatch(batcher, 0)
	if err != nil {
		t.Errorf("Linuxboot2uroot: %v", testsuite.DescribeBatcherErr(batcher, res, err))

	}

	err = r.PowerDown()
	if err != nil {
		t.Error(err)
	}

	err = r.Close()
	if err != nil {
		t.Error(err)
	}
}
