// Copyright 2020 Splitted-Desktop Systems. All rights reserved
// Copyright 2020 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
)

func TestSerialAndIPMI(t *testing.T) {
	ci, err := NewBMCControlledCi(true)
	defer CloseBMCControlledCiTest(t, ci)
	if err != nil {
		if _, ok := err.(SkipError); ok {
			t.Skipf("skipped (%v)", err)
		}
		t.Fatal(err)
	}

	err = ci.Open()
	if err != nil {
		t.Fatal(err)
	}

	batcher := testsuite.Linuxboot2urootBatcher
	// This Run will not return until the parallel tests finish.
	t.Run("group", func(t *testing.T) {
		t.Run("Serial", func(t *testing.T) {
			opts_s, warn := testsuite.ExpectOptions("TestSerialAndIPMI_serial")
			if warn != nil {
				t.Log(warn)
			}
			// spawn serial
			e, _, err := ci.Serial.Spawn(1*time.Second, opts_s...)
			if err != nil {
				t.Fatalf("Serial Spawn failed: %v", err)
			}
			out, _, err := e.Expect(regexp.MustCompile("LinuxBoot: Starting bzImage"), 120*time.Second)
			if err != nil {
				t.Errorf("error waiting for linuxboot loader: %v (got %v)", err, out)
			}

			t.Parallel()
			res, err := e.ExpectBatch(batcher, 0)
			if err != nil {
				t.Errorf("Linuxboot2uroot: %v", testsuite.DescribeBatcherErr(batcher, res, err))
			}
		})
		t.Run("IPMI", func(t *testing.T) {

			opts_i, warn := testsuite.ExpectOptions("TestSerialAndIPMI_ipmi")
			if warn != nil {
				t.Log(warn)
			}
			t.Parallel()

			// spawn ipmi
			e, _, err := ci.Ipmi.Spawn(1*time.Second, opts_i...)
			if err != nil {
				t.Fatalf("Spawn failed: %v", err)
			}
			res, err := e.ExpectBatch(batcher, 0)
			if err != nil {
				t.Errorf("Linuxboot2uroot: %v", testsuite.DescribeBatcherErr(batcher, res, err))
			}

			err = e.Close()
			if err != nil {
				t.Errorf("error closing expect: %v", err)
			}
		})
	}) // end of group

	t.Run("ipmi power", func(t *testing.T) {

		// Test IMPI Power interface
		on, err := ci.Ipmi.PowerStatus()
		if err != nil {
			t.Error(err)
		}
		if !on {
			t.Errorf("Expected on: PowerStatus()=%v, want %v", on, true)
		}
		err = ci.Ipmi.PowerDown()
		if err != nil {
			t.Error(err)
		}
		time.Sleep(1 * time.Second)
		on, err = ci.Ipmi.PowerStatus()
		if err != nil {
			t.Error(err)
		}
		if on {
			t.Errorf("Expected off: PowerStatus()=%v, want %v", on, false)
		}
	}) // end of ipmi power
}
