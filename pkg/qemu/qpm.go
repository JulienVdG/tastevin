// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package qemu

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/digitalocean/go-qemu/qmp"
	qmpraw "github.com/digitalocean/go-qemu/qmp/raw"
)

func (vm *VM) qpmOpen() error {
	var err error
	qmpsocket := filepath.Join(vm.tmpdir, qmpsock)
	vm.qmpmon, err = qmp.NewSocketMonitor("unix", qmpsocket, 2*time.Second)
	if err != nil {
		return fmt.Errorf("error creating qmp monitor: %v", err)
	}
	err = vm.qmpmon.Connect()
	if err != nil {
		return fmt.Errorf("error connecting qemu qmp monitor: %v", err)
	}

	vm.qmprawmon = qmpraw.NewMonitor(vm.qmpmon)

	si, err := vm.qmprawmon.QueryStatus()
	if err != nil {
		return fmt.Errorf("error querying status: %v", err)
	}
	if si.Running || si.Singlestep || si.Status != qmpraw.RunStatePrelaunch {
		return fmt.Errorf("wrong qemu starting state, expect {false false prelaunch}, got %v", si)
	}

	return nil
}
