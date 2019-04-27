// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package qemu implement the XXX interface on qemu
package qemu

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/digitalocean/go-qemu/qmp"
	qmpraw "github.com/digitalocean/go-qemu/qmp/raw"
	exp "github.com/google/goexpect"
)

const (
	qmpsock = "qmp-sock"
)

// Default configuration values
var (
	DefaultPath = "qemu-system-x86_64"
	DefaultArgs = []string{"-m", "1024",
		"-enable-kvm",
		"-object", "rng-random,filename=/dev/urandom,id=rng0",
		"-device", "virtio-rng-pci,rng=rng0",
	}
)

// VM represent a managed qemu session
type VM struct {
	exec.Cmd
	tmpdir     string
	cmdstarted bool
	qmpmon     *qmp.SocketMonitor
	qmprawmon  *qmpraw.Monitor
	expect     *exp.GExpect
	expectCh   <-chan error
}

// NewVM creates a new qemu vm
func NewVM(path string, extraArgs ...string) (*VM, error) {
	if path == "" {
		path = DefaultPath
	}
	args := append(DefaultArgs, extraArgs...)
	cmd := exec.Command(path, args...)
	return &VM{Cmd: *cmd}, nil
}

// Open a new VM session
func (vm *VM) Open() error {
	var err error
	vm.tmpdir, err = ioutil.TempDir("", "tastevin-qemu")
	if err != nil {
		return fmt.Errorf("error creating temp dir: %v", err)
	}
	qmpsocket := filepath.Join(vm.tmpdir, qmpsock)
	args := append(vm.Args, "-nographic", "-S", "-no-shutdown")
	args = append(args, "-qmp", fmt.Sprintf("unix:%s,server,nowait", qmpsocket))
	vm.Args = args

	err = vm.Start()
	if err != nil {
		vm.cleanup()
		return fmt.Errorf("error starting qemu: %v", err)
	}
	vm.cmdstarted = true
	// Wait for unix socket. TODO: timeout and parse errors
	timeout := time.After(5 * time.Second)
	tick := time.Tick(50 * time.Millisecond)
waitqmpsock:
	for {
		select {
		case <-timeout:
			return errors.New("timed out waiting for qmp server")
		case <-tick:
			_, err := os.Stat(qmpsocket)
			if err == nil {
				break waitqmpsock
			}
		}
	}
	err = vm.qpmOpen()
	if err != nil {
		vm.cleanup()
		return fmt.Errorf("error creating qmp monitor: %v", err)
	}

	return nil
}

// Close the VM session
func (vm *VM) Close() error {
	return vm.cleanup()
}

func (vm *VM) cleanup() error {
	var err error
	if vm.tmpdir != "" {
		os.RemoveAll(vm.tmpdir)
	}

	if vm.cmdstarted {
		if vm.qmprawmon != nil {
			err = vm.qmprawmon.Quit()
			if err != nil {
				return fmt.Errorf("error sending quit to qemu: %v", err)
			}
		} else {
			err = vm.Process.Kill()
			if err != nil {
				return fmt.Errorf("error killing qemu: %v", err)
			}
		}

		err = vm.Wait()
		if err != nil {
			return fmt.Errorf("error stopping qemu: %v", err)
		}
	}
	if vm.expect != nil {
		// try closing properly
		if vm.qmprawmon != nil {
			err = vm.qmprawmon.Quit()
			if err != nil {
				return fmt.Errorf("error sending quit to qemu: %v", err)
			}
		}
		err = vm.expect.Close()
		if err != nil {
			return fmt.Errorf("error closing expect: %v", err)
		}

		// Ensure the logs are closed by waiting the complete end
		<-vm.expectCh
	}

	if vm.qmpmon != nil {
		err = vm.qmpmon.Disconnect()
		if err != nil {
			return fmt.Errorf("error disconnecting monitor: %v", err)
		}
	}

	return nil
}

// PowerStatus returns the current power status
func (vm *VM) PowerStatus() (bool, error) {
	si, err := vm.qmprawmon.QueryStatus()
	if err != nil {
		return false, err
	}
	return si.Running, nil
}

// PowerUp starts the VM
func (vm *VM) PowerUp() error {
	err := vm.qmprawmon.SystemReset()
	if err != nil {
		return err
	}
	err = vm.qmprawmon.Cont()
	if err != nil {
		return err
	}
	return nil
}

// PowerDown stops the VM
func (vm *VM) PowerDown() error {
	err := vm.qmprawmon.Stop()
	if err != nil {
		return err
	}
	return nil
}
