// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
	exp "github.com/google/goexpect"
)

func TestBoot2(t *testing.T) {
	ci, err := NewWinterfellCi(false)
	defer CloseWinterfellCiTest(t, ci)
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

	bootBatcher := []exp.Batcher{
		&exp.BSnd{S: "boot2\r\n"},
		&testsuite.BExpTLog{
			L: "kexec done",
			R: "kexec_core: Starting new kernel",
			T: 5,
		},
		&testsuite.BExpTLog{
			L: "Debian booted",
			R: "debian-linuxboot-dut login:",
			T: 20,
		},
		&exp.BSnd{S: "root\r\n"},
		&exp.BExpT{R: "Password:", T: 1},
		&exp.BSnd{S: "r\r\n"},
		&testsuite.BExpTLog{
			L: "Logged in",
			R: "root@debian-linuxboot-dut:~#",
			T: 1,
		},
		&exp.BSnd{S: "dmesg\r\n"},
		&exp.BExpT{R: "Linux version", T: 1},
		&exp.BExpT{R: "root@debian-linuxboot-dut:~#", T: 20},
		&exp.BSnd{S: "poweroff\r\n"},
		&exp.BExpT{R: "reboot: Power down", T: 5},
	}

	res, err = e.ExpectBatch(bootBatcher, 0)
	if err != nil {
		t.Fatalf("Boot: %v", testsuite.DescribeBatcherErr(bootBatcher, res, err))
	}

}
