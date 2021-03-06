// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package em100 provides a Go API for controlling the em100 flash emulator
package em100

import (
	"log"
	"os"
	"os/exec"
)

// Default configuration values
var (
	DefaultPath = "/opt/em100/em100"
	DefaultChip = "W25Q128BV"
)

// Em100 is a flash emulator instance
type Em100 struct {
	args []string
}

// NewEm100 created a new Em100
func NewEm100(path, chip string) *Em100 {
	if path == "" {
		path = DefaultPath
	}
	if chip == "" {
		chip = DefaultChip
	}
	args := []string{path, "--set", chip}
	return &Em100{args: args}
}

func (em *Em100) run(extraArgs ...string) error {
	args := append(em.args, "--stop")
	args = append(args, extraArgs...)
	cmd := exec.Command(args[0], args[1:]...)
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
	// Check if file exists, but only log, command could run with sudo
	// and have access to a file not visible to us.
	info, err := os.Stat(filename)
	if err != nil {
		log.Println("cannot stat", filename, err)
	} else if info.IsDir() {
		log.Printf("error %s is a directory", filename)
	}
	return em.run("-d", filename, "--start")
}

// Dump a the emulator flash content to a file, this will stop the emulation
func (em *Em100) Dump(filename string) error {
	return em.run("-u", filename)
}
