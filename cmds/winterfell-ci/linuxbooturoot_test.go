// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/em100"
	"github.com/JulienVdG/tastevin/pkg/relay"
	"github.com/JulienVdG/tastevin/pkg/scriptreplay"
	"github.com/JulienVdG/tastevin/pkg/serial"
	"github.com/JulienVdG/tastevin/pkg/testsuite"
	exp "github.com/google/goexpect"
)

func TestLinuxboot2uroot(t *testing.T) {
	logdir := filepath.Join("testdata", "log")
	err := os.MkdirAll(logdir, 0775)
	if err != nil {
		t.Fatalf("TeeReplay failed: %v", err)
	}

	sr, err := scriptreplay.NewFileWriter(filepath.Join(logdir, "Linuxboot2uroot.log"), filepath.Join(logdir, "Linuxboot2uroot.tim"))
	if err != nil {
		t.Fatalf("TeeReplay failed: %v", err)
	}

	em := em100.NewEm100("", "")
	r := relay.NewRelay()

	// open serial
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600 /*, ReadTimeout: time.Nanosecond /*time.Second * 1.0 / 5760000*/}
	s, err := serial.NewSerial(c)
	if err != nil {
		log.Fatal(err)
	}

	em.Load("linuxboot.rom")

	err = r.PowerUp()
	if err != nil {
		t.Error(err)
	}

	// spawn serial
	e, _, err := s.Spawn(1*time.Second, exp.PartialMatch(true), exp.Tee(sr) /* exp.DebugCheck(nil), exp.Verbose(true)*/)
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

	err = sr.Close()
	if err != nil {
		t.Errorf("sr close: %v", err)
	}

}