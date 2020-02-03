// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
	exp "github.com/google/goexpect"
)

func TestLinuxboot2uroot(t *testing.T) {
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

	opts, warn := testsuite.ExpectOptions("")
	if warn != nil {
		t.Log(warn)
	}

	// spawn serial
	e, _, err := ci.Serial.Spawn(1*time.Second, opts...)
	if err != nil {
		t.Fatalf("Serial Spawn failed: %v", err)
	}

	batcher := testsuite.LinuxbootEfiLoaderBatcher
	batcher = append(batcher, testsuite.Linuxboot2urootBatcher...)
	res, err := e.ExpectBatch(batcher, 0)
	if err != nil {
		t.Fatalf("Linuxboot2uroot: %v", testsuite.DescribeBatcherErr(batcher, res, err))
	}
	if t.Failed() {
		fmt.Printf("Initial boot fail, try cold reboot without reloading the flash...\n")
		success := t.Run("ColdReboot", func(t *testing.T) {
			err := ci.Ipmi.PowerDown()
			if err != nil {
				t.Fatalf("cannot power down: %v", err)
			}
			err = ci.Ipmi.PowerUp()
			if err != nil {
				t.Fatalf("cannot power up: %v", err)
			}
			fmt.Printf("Cold Reboot done\n")

			res, err := e.ExpectBatch(batcher, 0)
			if err != nil {
				t.Fatalf("Linuxboot2uroot: %v", testsuite.DescribeBatcherErr(batcher, res, err))
			}
		})

		if success {
			// printed in subtest: https://github.com/golang/go/issues/29755
			// + https://github.com/golang/go/issues/24929
			t.Log("Initial boot fail, cold reboot without reloading the flash succeeded. Something is wrong with the initial state of the flash!\n")
		} else {
			t.FailNow()
		}
	}

	t.Run("Reboot", func(t *testing.T) {
		rebootBatcher := []exp.Batcher{
			&exp.BSnd{S: "cat >proc/sysrq-trigger\r\n"},
			&exp.BExpT{R: "sysrq-trigger", T: 1},
			&exp.BSnd{S: "b\r\n"},
			&testsuite.BExpTLog{
				L: "Reboot done",
				R: "sysrq:( SysRq :)? Resetting",
				T: 1,
			}}

		res, err := e.ExpectBatch(rebootBatcher, 0)
		if err != nil {
			t.Fatalf("Reboot: %v", testsuite.DescribeBatcherErr(rebootBatcher, res, err))
		}

		res, err = e.ExpectBatch(batcher, 0)
		if err != nil {
			t.Fatalf("Linuxboot2uroot: %v", testsuite.DescribeBatcherErr(batcher, res, err))
		}
	})
}
