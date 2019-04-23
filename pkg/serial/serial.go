// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package serial implement the XXX interface over serial
package serial

import (
	serialport "github.com/tarm/serial"
)

// Config contains the information needed to open a serial port.
type Config serialport.Config

// Serial represent a serial port session
type Serial struct {
	c Config
	p *serialport.Port
}

// NewSerial creates a new serial port session
func NewSerial(c *Config) (*Serial, error) {
	return &Serial{c: *c}, nil
}

// Open the serial session
func (s *Serial) Open() error {
	var err error
	c := serialport.Config(s.c)
	s.p, err = serialport.OpenPort(&c)
	return err
}

// Close the serial session
func (s *Serial) Close() error {
	return s.p.Close()

}
