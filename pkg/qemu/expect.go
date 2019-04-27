// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qemu

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"

	exp "github.com/google/goexpect"
)

// Spawn return a goexpect.GExpect for the VM
// TODO: Open should not be called before Spawn for qemu while it must for ipmi
func (vm *VM) Spawn(timeout time.Duration, opts ...exp.Option) (*exp.GExpect, <-chan error, error) {
	var err error
	vm.tmpdir, err = ioutil.TempDir("", "tastevin-qemu")
	if err != nil {
		return nil, nil, fmt.Errorf("error creating temp dir: %v", err)
	}
	qmpsocket := filepath.Join(vm.tmpdir, qmpsock)
	args := append(vm.Args, "-nographic", "-S", "-no-shutdown")
	args = append(args, "-qmp", fmt.Sprintf("unix:%s,server", qmpsocket))
	vm.Args = args

	e, errchan, err := exp.SpawnWithArgs(vm.Args, timeout, opts...)
	if err != nil {
		vm.cleanup()
		return nil, nil, fmt.Errorf("error spawning qemu: %v", err)
	}
	vm.expect = e
	vm.expectCh = errchan

	// qemu-system-x86_64: -qmp unix:/tmp/tastevin-qemu238249216/qmp-sock,server: info: QEMU waiting for connection on: disconnected:unix:/tmp/tastevin-qemu238249216/qmp-sock,server
	out, _, err := e.Expect(regexp.MustCompile("-qmp .* info: QEMU waiting for connection on: disconnected:"), 1*time.Second)
	if err != nil {
		vm.cleanup()
		return nil, nil, fmt.Errorf("error waiting for qemu qmp sever: %v (got %v)", err, out)
	}

	err = vm.qpmOpen()
	if err != nil {
		vm.cleanup()
		return nil, nil, fmt.Errorf("error creating qmp monitor: %v", err)
	}

	return e, errchan, err
}
