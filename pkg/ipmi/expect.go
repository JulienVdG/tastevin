// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	exp "github.com/google/goexpect"
)

// options is a copy of goipmi.tool private function
func (r *Remote) options() []string {
	intf := r.Interface
	if intf == "" {
		intf = "lanplus"
	}

	options := []string{
		"-H", r.Hostname,
		"-U", r.Username,
		"-P", r.Password,
		"-I", intf,
	}

	if r.Port != 0 {
		options = append(options, "-p", strconv.Itoa(r.Port))
	}

	return options
}

// cmdSOLArgs returns the commandline similar to goipmi.tool Console
func (r *Remote) cmdSOLArgs() []string {
	path := r.Path

	if path == "" {
		path = "ipmitool"
	}

	args := []string{path}
	args = append(args, r.options()...)
	args = append(args, "sol", "activate")

	return args
}

// Spawn creates a goexpect session over IPMI SOL
func (r *Remote) Spawn(timeout time.Duration, opts ...exp.Option) (*exp.GExpect, <-chan error, error) {
	// TODO: Should also call open
	e, errchan, err := exp.SpawnWithArgs(r.cmdSOLArgs(), timeout, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("error spawning impi sol: %v", err)
	}
	// TODO match first '[SOL Session operational.  Use ~? for help]'
	out, _, err := e.Expect(regexp.MustCompile("\\[SOL Session operational.  Use .\\? for help\\]"), 20*time.Second)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening ipmi sol: %v (got %v)", err, out)
	}

	// TODO close should probably deactivate or TERM the process

	return e, errchan, err
}
