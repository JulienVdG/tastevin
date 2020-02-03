// Copyright 2020 Splitted-Desktop Systems. All rights reserved
// Copyright 2020 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import "os"

// GetLinuxBootImage returns the path to the bios image
// from TASTEVIN_LINUXBOOT_IMAGE env var, default to "linuxboot.rom"
// example env:
// TASTEVIN_LINUXBOOT_IMAGE='/home/linuxboot/linuxboot.rom'
func GetLinuxBootImage() string {
	env := os.Getenv("TASTEVIN_LINUXBOOT_IMAGE")
	if env == "" {
		return "linuxboot.rom"
	}
	return env
}
