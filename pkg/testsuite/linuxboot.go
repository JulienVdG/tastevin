// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testsuite

import (
	exp "github.com/google/goexpect"
)

// LinuxbootEfiLoaderBatcher follows the boot sequence of linuxboot efi loader
var LinuxbootEfiLoaderBatcher []exp.Batcher = []exp.Batcher{
	&BExpTLog{
		L: "Matched LinuxBoot banner",
		R: "\\| Starting LinuxBoot \\|",
		T: 60,
	}, &BExpTLog{
		L: "Matched Starting bzImage",
		R: "LinuxBoot: Starting bzImage",
		T: 30,
	}}

// Linuxboot2urootBatcher follows the boot sequence of u-root to the shell prompt
var Linuxboot2urootBatcher []exp.Batcher = []exp.Batcher{
	&BExpTLog{
		L: "Matched u-root banner",
		R: "Welcome to u-root!",
		T: 40, // TODO make this time a parameter
	}, &BExpTLog{
		L: "Matched u-root prompt",
		R: "~/> ",
		T: 20,
	}}
