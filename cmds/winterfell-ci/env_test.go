// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import "os"

// GetWinterfellImage returns the path to the bios image
// from TASTEVIN_WINTERFELL_IMAGE env var, default to "linuxboot.rom"
// example env:
// TASTEVIN_WINTERFELL_IMAGE='/home/linuxboot/linuxboot.rom'
func GetWinterfellImage() string {
	env := os.Getenv("TASTEVIN_WINTERFELL_IMAGE")
	if env == "" {
		return "linuxboot.rom"
	}
	return env
}
