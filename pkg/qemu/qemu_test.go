// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qemu

import (
	"os"
	"testing"
	"time"
)

func TestQemu(t *testing.T) {
	vm, err := NewVM("",
		"-kernel", "/boot/vmlinuz-4.19.0-2-amd64",
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
		t.Errorf("Expected off got %v", err)
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
		t.Errorf("Expected on got %v", err)
	}

	time.Sleep(3 * time.Second)

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
		t.Errorf("Expected off got %v", err)
	}

	err = vm.Close()
	if err != nil {
		t.Error(err)
	}

}
