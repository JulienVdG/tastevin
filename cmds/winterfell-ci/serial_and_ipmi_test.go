// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/em100"
	"github.com/JulienVdG/tastevin/pkg/ipmi"
	"github.com/JulienVdG/tastevin/pkg/relay"
	"github.com/JulienVdG/tastevin/pkg/serial"
	"github.com/JulienVdG/tastevin/pkg/testsuite"
)

func TestSerialAndIPMI(t *testing.T) {
	em, err := em100.NewEm100FromEnv()
	if err != nil {
		t.Skipf("skipped unless TASTEVIN_EM100 is set (%v)", err)
	}
	r, err := relay.NewRelayFromEnv()
	if err != nil {
		t.Skipf("skipped unless TASTEVIN_RELAY is set (%v)", err)
	}

	ic, err := ipmi.ConfigFromEnv()
	if err != nil {
		t.Skipf("IPMI test is skipped unless TASTEVIN_IPMI is set (%v)", err)
	}
	i, err := ipmi.NewRemote(ic)
	if err != nil {
		t.Fatal(err)
	}

	opts_s, warn := testsuite.ExpectOptions("TestSerialAndIPMI_serial")
	if warn != nil {
		t.Log(warn)
	}

	opts_i, warn := testsuite.ExpectOptions("TestSerialAndIPMI_ipmi")
	if warn != nil {
		t.Log(warn)
	}

	// open serial
	sc := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600 /*, ReadTimeout: time.Nanosecond /*time.Second * 1.0 / 5760000*/}
	s, err := serial.NewSerial(sc)
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

	err = i.Open()
	if err != nil {
		t.Error(err)
	}

	// This Run will not return until the parallel tests finish.
	t.Run("group", func(t *testing.T) {
		t.Run("Serial", func(t *testing.T) {
			// spawn serial
			e, _, err := s.Spawn(1*time.Second, opts_s...)
			if err != nil {
				t.Fatalf("Serial Spawn failed: %v", err)
			}
			out, _, err := e.Expect(regexp.MustCompile("LinuxBoot: Starting bzImage"), 30*time.Second)
			if err != nil {
				t.Errorf("error waiting for linuxboot loader: %v (got %v)", err, out)
			}

			t.Parallel()
			err = testsuite.Linuxboot2uroot(t, e)
			if err != nil {
				t.Errorf("Linuxboot2uroot returned: %v", err)
			}
		})
		t.Run("IPMI", func(t *testing.T) {
			t.Parallel()

			e, _, err := i.Spawn(1*time.Second, opts_i...)
			if err != nil {
				t.Fatalf("Spawn failed: %v", err)
			}
			err = testsuite.Linuxboot2uroot(t, e)
			if err != nil {
				t.Errorf("Linuxboot2uroot returned: %v", err)
			}

			err = e.Close()
			if err != nil {
				t.Errorf("error closing expect: %v", err)
			}
		})
	}) // end of group

	t.Run("ipmi power", func(t *testing.T) {

		// Test IMPI Power interface
		on, err := i.PowerStatus()
		if err != nil {
			t.Error(err)
		}
		if !on {
			t.Error("Expected Power status On, got Off")
		}
		err = i.PowerDown()
		if err != nil {
			t.Error(err)
		}
		time.Sleep(1 * time.Second)
		on, err = i.PowerStatus()
		if err != nil {
			t.Error(err)
		}
		if on {
			t.Error("Expected Power status Off, got On")
		}
	}) // end of ipmi power

	err = i.Close()
	if err != nil {
		t.Error(err)
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
