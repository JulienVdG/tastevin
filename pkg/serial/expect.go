// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serial

import (
	"fmt"
	"time"

	exp "github.com/google/goexpect"
)

// Spawn creates a goexpect session over serial
func (s *Serial) Spawn(timeout time.Duration, opts ...exp.Option) (*exp.GExpect, <-chan error, error) {
	var err error

	err = s.Open()
	if err != nil {
		return nil, nil, fmt.Errorf("error opening serial port: %v", err)
	}

	resCh := make(chan error)

	e, errchan, err := exp.SpawnGeneric(&exp.GenOptions{
		In:  s.p,
		Out: s.p,
		Wait: func() error {
			return <-resCh
		},
		Close: func() error {
			close(resCh)
			return s.Close()
		},
		Check: func() bool { return true },
	}, timeout, opts...)
	if err != nil {
		s.Close()
		return nil, nil, fmt.Errorf("error spawning serial port: %v", err)
	}

	return e, errchan, err
}
