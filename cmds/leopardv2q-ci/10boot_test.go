// Copyright 2020 Splitted-Desktop Systems. All rights reserved
// Copyright 2020 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
	exp "github.com/google/goexpect"
)

// elvish echoes the prompt after each received byte, also color the cmd :(
// send by chuck and expect echo (assume color change on space boundary)
func genSendExpectEchoLn(msg string) []exp.Batcher {
	var b []exp.Batcher
	const chunk = 16
	var end, endnosp, space int
	l := len(msg)

	for i := 0; i < l; i = end {
		end = i + chunk
		endnosp = end
		sp := i + strings.IndexAny(msg[i:], " \t") + 1
		if sp != i && sp < end {
			end = sp
			endnosp = sp - 1
		}
		if end > l {
			end = l
			endnosp = l
		}
		//fmt.Printf("i:%d,space:%d,end:%d,sp:%d,len:%d\n", i, space, end, sp, l)
		b = append(b,
			&exp.BSnd{S: msg[i:end]},
			&exp.BExpT{R: regexp.QuoteMeta(msg[space:endnosp]), T: 1},
		)
		if sp <= end {
			space = sp
		}
	}
	b = append(b, &exp.BSnd{S: "\r\n"})
	return b
}

func disabledTestGenSendExpectEchoLn(t *testing.T) {
	bootBatcher := genSendExpectEchoLn("boot -remove '$vt_handoff,console,quiet,splash' -reuse earlyprintk,console")
	fmt.Printf("%#q", bootBatcher)
}

func TestBoot(t *testing.T) {
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

	batcher := testsuite.Linuxboot2urootBatcher
	res, err := e.ExpectBatch(batcher, 0)
	if err != nil {
		t.Fatalf("Linuxboot2uroot: %v", testsuite.DescribeBatcherErr(batcher, res, err))
	}

	bootBatcher := genSendExpectEchoLn("boot -remove '$vt_handoff,console,quiet,splash' -reuse earlyprintk,console")

	bootBatcher = append(bootBatcher, []exp.Batcher{
		&testsuite.BExpTLog{
			L: "kexec done",
			R: "kexec_core: Starting new kernel",
			T: 10,
		},
		&testsuite.BExpTLog{
			L: "Ubuntu booted",
			R: "Ubuntu login:",
			T: 120,
		},
		&exp.BSnd{S: "ubuntu\r\n"},
		&exp.BExpT{R: "Password:", T: 1},
		&exp.BSnd{S: "ubuntu\r\n"},
		&testsuite.BExpTLog{
			L: "Logged in",
			R: "ubuntu@ubuntu:~$",
			T: 1,
		},
		&exp.BSnd{S: "dmesg\r\n"},
		&exp.BExpT{R: "Linux version", T: 1},
		&exp.BExpT{R: "ubuntu@ubuntu:~$", T: 20},
		&exp.BSnd{S: "poweroff\r\n"},
		&exp.BExpT{R: "reboot: Power down", T: 10},
	}...)

	res, err = e.ExpectBatch(bootBatcher, 0)
	if err != nil {
		t.Fatalf("Boot: %v", testsuite.DescribeBatcherErr(bootBatcher, res, err))
	}

}
