// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package relay provides a Go API for controlling the relays that powers nodes
package relay

import "os/exec"

// TODO implement in pure go (using a libFTDI go wrapper)
// See:
// https://learn.adafruit.com/adafruit-ft232h-breakout/mpsse-setup
// https://learn.adafruit.com/adafruit-ft232h-breakout/linux-setup
// https://www.intra2net.com/en/developer/libftdi/
// https://www.intra2net.com/en/developer/libftdi/documentation/
// https://github.com/adafruit/Adafruit_Python_GPIO/blob/master/Adafruit_GPIO/FT232H.py
// https://github.com/stvnrhodes/goftdi/blob/master/ftdi.go

// Hardcode script names as this is only a temporary solution cf TODO above.
// A solution that will last as they always do :)

const (
	powerOnCmd  = "/usr/local/bin/relay-poweron"
	powerOffCmd = "/usr/local/bin/relay-poweroff"
)

// Relay represent a relay controlling a node power
type Relay struct{}

// NewRelay creates a Relay
func NewRelay() *Relay {
	return &Relay{}
}

// TODO do we need Open/Close ?

// PowerUp starts the node
func (r *Relay) PowerUp() error {
	cmd := exec.Command(powerOnCmd)
	return cmd.Run()
}

// PowerDown stops the node
func (r *Relay) PowerDown() error {
	cmd := exec.Command(powerOffCmd)
	return cmd.Run()
}
