// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qemu

import (
	"os"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/testsuite"
)

func TestLinuxboot2uroot(t *testing.T) {
	if _, ok := os.LookupEnv("TASTEVIN_QEMU"); !ok {
		t.Skip("QEMU test is skipped unless TASTEVIN_QEMU is set")
	}

	vm, err := NewVM("",
		"-kernel", "/boot/vmlinuz-4.19.0-2-amd64",
		"-initrd", "/tmp/initramfs.linux_amd64.cpio",
		"-append", "console=ttyS0")
	if err != nil {
		t.Fatal(err)
	}

	opts, warn := testsuite.ExpectOptions("")
	if warn != nil {
		t.Log(warn)
	}

	e, _, err := vm.Spawn(1*time.Second, opts...)
	if err != nil {
		t.Fatalf("Spawn failed: %v", err)
	}
	err = vm.PowerUp()
	if err != nil {
		t.Error(err)
	}

	err = testsuite.Linuxboot2uroot(t, e)
	if err != nil {
		t.Errorf("Linuxboot2uroot returned: %v", err)
	}

	err = vm.PowerDown()
	if err != nil {
		t.Error(err)
	}

	err = vm.Close()
	if err != nil {
		t.Error(err)
	}
}
