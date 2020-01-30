// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

// Load a flash image into the BIOS flash using the BMC
func (r *Remote) Load(filename string) error {
	if r.YafuFlashPath == "" {
		return errors.New("YafuFlashPath not configured")
	}
	if s, err := r.PowerStatus(); err != nil {
		return err
	} else if s {
		return errors.New("Cannot update BIOS flash while system is running")
	}
	// Check if file exists, but only log, command could run with sudo
	// and have access to a file not visible to us.
	info, err := os.Stat(filename)
	if err != nil {
		log.Println("cannot stat", filename, err)
	} else if info.IsDir() {
		log.Printf("error %s is a directory", filename)
	}
	return r.yafuflash(filename)
}

func (r *Remote) yafuflash(filename string) error {
	if r.YafuFlashPath == "" {
		return errors.New("YafuFlashPath not configured")
	}
	args := []string{"-nw", "-u", r.Username, "-p", r.Password}
	if net.ParseIP(r.Hostname) == nil {
		args = append(args, "-host", r.Hostname)
	} else {
		args = append(args, "-ip", r.Hostname)
	}
	if r.Port != 0 {
		args = append(args, "-P", fmt.Sprintf("%d", r.Port))
	}
	args = append(args, "-bios", filename)

	// TODO: do we want to use expect, monitor the output and handle errors
	cmd := exec.Command(r.YafuFlashPath, args...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running YafuFlash: %v", err)
	}
	return nil
}
