// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package serial implement the XXX interface over serial
package serial

import (
	"fmt"
	"io"
	"time"

	exp "github.com/google/goexpect"
	serialport "github.com/tarm/serial"
)

// Config contains the information needed to open a serial port.
type Config serialport.Config

// Serial represent a serial port session
type Serial struct {
	c        Config
	p        *serialport.Port
	expect   *exp.GExpect
	expectCh <-chan error
}

// NewSerial creates a new serial port session
func NewSerial(c *Config) (*Serial, error) {
	return &Serial{c: *c}, nil
}

// Open the serial session
func (s *Serial) Open() error {
	var err error
	c := serialport.Config(s.c)
	if c.ReadTimeout == 0 {
		c.ReadTimeout = time.Nanosecond
	}
	s.p, err = serialport.OpenPort(&c)
	return err
}

// Close the serial session
func (s *Serial) Close() error {
	if s.p == nil {
		return nil
	}
	if s.expect != nil {
		err := s.expect.Close()
		if err != nil {
			return fmt.Errorf("error closing expect: %v", err)
		}
		// Ensure the logs are closed by waiting the complete end
		<-s.expectCh
		// expect will close the serial port, exit now
		return nil
	}
	err := s.p.Close()
	s.p = nil
	return err
}

func (s *Serial) Read(b []byte) (n int, err error) {
	if s.p == nil {
		return 0, fmt.Errorf("read %s port closed", s.c.Name)
	}
	n, err = s.p.Read(b)
	// Handle read timeout, should not return EOF until closed
	if n == 0 && err == io.EOF && s.p != nil {
		err = nil
	}
	return
}

func (s *Serial) Write(b []byte) (n int, err error) {
	if s.p == nil {
		return 0, fmt.Errorf("write %s port closed", s.c.Name)
	}
	return s.p.Write(b)
}
