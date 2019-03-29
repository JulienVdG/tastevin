// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qemu

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/scriptreplay"
	"github.com/JulienVdG/tastevin/pkg/testsuite"
	exp "github.com/google/goexpect"
)

func TestLinuxboot2uroot(t *testing.T) {
	vm, err := NewVM("",
		"-kernel", "/boot/vmlinuz-4.19.0-2-amd64",
		"-initrd", "/tmp/initramfs.linux_amd64.cpio",
		"-append", "console=ttyS0")
	if err != nil {
		t.Fatal(err)
	}

	logdir := filepath.Join("testdata", "log")
	err = os.MkdirAll(logdir, 0775)
	if err != nil {
		t.Fatalf("TeeReplay failed: %v", err)
	}

	sr, err := scriptreplay.NewFileWriter(filepath.Join(logdir, "QemuLinuxboot2uroot.log"), filepath.Join(logdir, "QemuLinuxboot2uroot.tim"))
	if err != nil {
		t.Fatalf("TeeReplay failed: %v", err)
	}
	e, _, err := vm.Spawn(1*time.Second, exp.PartialMatch(true), exp.Tee(sr) /* exp.DebugCheck(nil), exp.Verbose(true)*/)
	if err != nil {
		t.Fatalf("Spawn failed: %v", err)
	}
	err = vm.PowerUp()
	if err != nil {
		t.Error(err)
	}

	err = testsuite.Linuxboot2uroot(t, e)
	if err != nil {
		t.Fatalf("Linuxboot2uroot returned: %v", err)
	}

	err = vm.PowerDown()
	if err != nil {
		t.Error(err)
	}

	err = vm.Close()
	if err != nil {
		t.Error(err)
	}

	err = sr.Close()
	if err != nil {
		t.Errorf("sr close: %v", err)
	}
}
