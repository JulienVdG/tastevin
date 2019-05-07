// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/em100"
	"github.com/JulienVdG/tastevin/pkg/relay"
	"github.com/JulienVdG/tastevin/pkg/serial"
	"github.com/JulienVdG/tastevin/pkg/testsuite"
)

func TestLinuxboot2uroot(t *testing.T) {
	em, err := em100.NewEm100FromEnv()
	if err != nil {
		t.Skipf("skipped unless TASTEVIN_EM100 is set (%v)", err)
	}
	r, err := relay.NewRelayFromEnv()
	if err != nil {
		t.Skipf("skipped unless TASTEVIN_RELAY is set (%v)", err)
	}

	opts, warn := testsuite.ExpectOptions("")
	if warn != nil {
		t.Log(warn)
	}

	// open serial
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600 /*, ReadTimeout: time.Nanosecond /*time.Second * 1.0 / 5760000*/}
	s, err := serial.NewSerial(c)
	if err != nil {
		t.Fatal(err)
	}
	err = s.Open()
	if err != nil {
		t.Fatal(err)
	}

	em.Load(GetWinterfellImage())

	err = r.PowerUp()
	if err != nil {
		t.Error(err)
	}

	// spawn serial
	e, _, err := s.Spawn(1*time.Second, opts...)
	if err != nil {
		t.Fatalf("Serial Spawn failed: %v", err)
	}

	err = testsuite.Linuxboot2uroot(t, e)
	if err != nil {
		t.Fatalf("Linuxboot2uroot returned: %v", err)
	}

	err = r.PowerDown()
	if err != nil {
		t.Error(err)
	}

	em.Stop()

	err = s.Close()
	if err != nil {
		t.Error(err)
	}
}
