// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package em100 provides a Go API for controlling the em100 flash emulator
package em100

import "os/exec"

// Default configuration values
var (
	DefaultPath = "/opt/em100/em100"
	DefaultChip = "W25Q128BV"
)

// Em100 is a flash emulator instance
type Em100 struct {
	path string
	chip string
}

// NewEm100 created a new Em100
func NewEm100(path, chip string) *Em100 {
	if path == "" {
		path = DefaultPath
	}
	if chip == "" {
		chip = DefaultChip
	}
	return &Em100{path: path, chip: chip}
}

func (em *Em100) run(extraArgs ...string) error {
	args := []string{"--stop", "--set", em.chip}
	args = append(args, extraArgs...)
	cmd := exec.Command(em.path, args...)
	return cmd.Run()
}

// Start the emulation
func (em *Em100) Start() error {
	return em.run("--start")
}

// Stop the emulation
func (em *Em100) Stop() error {
	return em.run()
}

// Load a flash image into the emulator, this will start the emulation
func (em *Em100) Load(filename string) error {
	return em.run("-d", filename, "--start")
}

// Dump a the emulator flash content to a file, this will stop the emulation
func (em *Em100) Dump(filename string) error {
	return em.run("-u", filename)
}
