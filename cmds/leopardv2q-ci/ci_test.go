// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/JulienVdG/tastevin/pkg/ipmi"
	"github.com/JulienVdG/tastevin/pkg/serial"
)

type BMCControlledCi struct {
	Serial *serial.Serial
	Ipmi   *ipmi.Remote
}

type SkipError error

func NewBMCControlledCi(withSerial bool) (*BMCControlledCi, error) {
	ci := &BMCControlledCi{}
	var err error

	// check ipmi config first (skip case)
	ic, err := ipmi.ConfigFromEnv()
	if err != nil {
		return ci, SkipError(err)
	}
	ci.Ipmi, err = ipmi.NewRemote(ic)
	if err != nil {
		return ci, err
	}
	if withSerial {
		// create serial
		sc := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600}
		ci.Serial, err = serial.NewSerial(sc)
		if err != nil {
			return ci, err
		}
	}

	return ci, nil
}

func (ci *BMCControlledCi) Open() error {
	err := ci.Ipmi.Open()
	if err != nil {
		return err
	}

	err = ci.Ipmi.Load(GetLinuxBootImage())
	if err != nil {
		return err
	}

	if ci.Serial != nil {

		err = ci.Serial.Open()
		if err != nil {
			return err
		}
	}

	err = ci.Ipmi.PowerUp()
	if err != nil {
		return err
	}

	return nil
}

func (ci *BMCControlledCi) Close() error {
	var msg []string
	if ci.Ipmi != nil {
		err := ci.Ipmi.PowerDown()
		if err != nil {
			if err != nil {
				msg = append(msg, fmt.Sprintf("cannot power down with IPMI (err: %v)", err))
			}
		}
	}

	if ci.Serial != nil {
		err := ci.Serial.Close()

		if err != nil {
			msg = append(msg, fmt.Sprintf("cannot close serial (err: %v)", err))
		}
	}
	if ci.Ipmi != nil {
		err := ci.Ipmi.Close()

		if err != nil {
			msg = append(msg, fmt.Sprintf("cannot close IPMI (err: %v)", err))
		}
	}
	if len(msg) > 0 {
		return errors.New(strings.Join(msg, "; "))
	}
	return nil
}

func CloseBMCControlledCiTest(t *testing.T, ci *BMCControlledCi) {
	err := ci.Close()
	if err != nil {
		t.Fatal(err)
	}
}
