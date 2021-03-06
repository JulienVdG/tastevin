// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qemu

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestQemu(t *testing.T) {
	if _, ok := os.LookupEnv("TASTEVIN_QEMU"); !ok {
		t.Skip("QEMU test is skipped unless TASTEVIN_QEMU is set")
	}

	vm, err := NewVM("",
		"-kernel", "/boot/vmlinuz-4.19.0-5-amd64",
		"-initrd", "/tmp/initramfs.linux_amd64.cpio",
		"-append", "console=ttyS0")
	if err != nil {
		t.Fatal(err)
	}
	// test
	vm.Stdin = os.Stdin
	vm.Stdout = os.Stdout
	vm.Stderr = os.Stderr

	err = vm.Open()
	if err != nil {
		t.Fatal(err)
	}

	on, err := vm.PowerStatus()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Power status is %v", on)
	if on {
		t.Errorf("Expected off: PowerStatus()=%v, want %v", on, false)
	}

	err = vm.PowerUp()
	if err != nil {
		t.Error(err)
	}
	on, err = vm.PowerStatus()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Power status is %v", on)
	if !on {
		t.Errorf("Expected on: PowerStatus()=%v, want %v", on, true)
	}

	time.Sleep(10 * time.Second)

	err = vm.PowerDown()
	if err != nil {
		t.Error(err)
	}
	on, err = vm.PowerStatus()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Power status is %v", on)
	if on {
		t.Errorf("Expected off: PowerStatus()=%v, want %v", on, false)
	}

	err = vm.Close()
	if err != nil {
		t.Error(err)
	}

	// We mess with stdout, add a newline
	fmt.Println("")
}
