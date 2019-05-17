// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
)

func TestLinuxboot2uroot(t *testing.T) {
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

	err = testsuite.Linuxboot2uroot(t, e)
	if err != nil {
		t.Fatalf("Linuxboot2uroot returned: %v", err)
	}

	t.Run("Reboot", func(t *testing.T) {
		err := e.Send("cat >proc/sysrq-trigger\r\n")
		if err != nil {
			t.Fatalf("Open sysrq: %v", err)
		}
		out, _, err := e.Expect(regexp.MustCompile("sysrq-trigger"), 1*time.Second)
		if err != nil {
			t.Errorf("error waiting for sysrq open: %v (got %v)", err, out)
		}

		err = e.Send("b\r\n")
		if err != nil {
			t.Fatalf("Rebooting: %v", err)
		}
		out, _, err = e.Expect(regexp.MustCompile("sysrq: SysRq : Resetting"), 1*time.Second)
		if err != nil {
			t.Errorf("error waiting for sysrq reset: %v (got %v)", err, out)
		}
		fmt.Printf("Reboot done\n")

		out, _, err = e.Expect(regexp.MustCompile("LinuxBoot: Starting bzImage"), 30*time.Second)
		if err != nil {
			t.Errorf("error waiting for linuxboot loader: %v (got %v)", err, out)
		}

		err = testsuite.Linuxboot2uroot(t, e)
		if err != nil {
			t.Fatalf("Linuxboot2uroot returned: %v", err)
		}
	})
}
