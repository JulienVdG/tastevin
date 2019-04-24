// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/em100"
	"github.com/JulienVdG/tastevin/pkg/ipmi"
	"github.com/JulienVdG/tastevin/pkg/relay"
	"github.com/JulienVdG/tastevin/pkg/scriptreplay"
	"github.com/JulienVdG/tastevin/pkg/serial"
	"github.com/JulienVdG/tastevin/pkg/testsuite"
	exp "github.com/google/goexpect"
)

func TestSerialAndIPMI(t *testing.T) {
	ic, err := ipmi.ConfigFromEnv()
	if err != nil {
		t.Skipf("IPMI test is skipped unless TASTEVIN_IPMI is set (%v)", err)
	}
	i, err := ipmi.NewRemote(ic)
	if err != nil {
		t.Fatal(err)
	}

	logdir := filepath.Join("testdata", "log")
	err = os.MkdirAll(logdir, 0775)
	if err != nil {
		t.Fatalf("TeeReplay failed: %v", err)
	}

	sr_s, err := scriptreplay.NewFileWriter(filepath.Join(logdir, "TestSerialAndIPMI_serial.log"), filepath.Join(logdir, "TestSerialAndIPMI_serial.tim"))
	if err != nil {
		t.Fatalf("TeeReplay failed: %v", err)
	}

	sr_i, err := scriptreplay.NewFileWriter(filepath.Join(logdir, "TestSerialAndIPMI_ipmi.log"), filepath.Join(logdir, "TestSerialAndIPMI_ipmi.tim"))
	if err != nil {
		t.Fatalf("TeeReplay failed: %v", err)
	}

	em := em100.NewEm100("", "")
	r := relay.NewRelay()

	// open serial
	sc := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600 /*, ReadTimeout: time.Nanosecond /*time.Second * 1.0 / 5760000*/}
	s, err := serial.NewSerial(sc)
	if err != nil {
		log.Fatal(err)
	}

	em.Load("linuxboot.rom")

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
			e, _, err := s.Spawn(1*time.Second, exp.PartialMatch(true), exp.Tee(sr_s))
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

			err = sr_s.Close()
			if err != nil {
				t.Errorf("sr close: %v", err)
			}
		})
		t.Run("IPMI", func(t *testing.T) {
			t.Parallel()

			e, _, err := i.Spawn(1*time.Second, exp.PartialMatch(true), exp.Tee(sr_i))
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

			err = sr_i.Close()
			if err != nil {
				t.Errorf("sr close: %v", err)
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
			t.Error("Expected Power status On, go Off")
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
			t.Error("Expected Power status Off, go On")
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
