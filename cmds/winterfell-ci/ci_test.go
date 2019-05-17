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

	"github.com/JulienVdG/tastevin/pkg/em100"
	"github.com/JulienVdG/tastevin/pkg/ipmi"
	"github.com/JulienVdG/tastevin/pkg/relay"
	"github.com/JulienVdG/tastevin/pkg/serial"
)

type WinterfellCi struct {
	Em100  *em100.Em100
	Relay  *relay.Relay
	Serial *serial.Serial
	Ipmi   *ipmi.Remote
}

type SkipError error

func NewWinterfellCi(withIpmi bool) (*WinterfellCi, error) {
	ci := &WinterfellCi{}
	var err error
	ci.Em100, err = em100.NewEm100FromEnv()
	if err != nil {
		return ci, SkipError(err)
	}
	ci.Relay, err = relay.NewRelayFromEnv()
	if err != nil {
		return ci, SkipError(err)
	}
	// check ipmi config first (skip case)
	if withIpmi {
		ic, err := ipmi.ConfigFromEnv()
		if err != nil {
			return ci, SkipError(err)
		}
		ci.Ipmi, err = ipmi.NewRemote(ic)
		if err != nil {
			return ci, err
		}
	}
	// create serial
	sc := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600}
	ci.Serial, err = serial.NewSerial(sc)
	if err != nil {
		return ci, err
	}

	return ci, nil
}

func (ci *WinterfellCi) Open() error {
	err := ci.Serial.Open()
	if err != nil {
		return err
	}

	ci.Em100.Load(GetWinterfellImage())

	err = ci.Relay.PowerUp()
	if err != nil {
		return err
	}

	if ci.Ipmi != nil {

		err = ci.Ipmi.Open()
		if err != nil {
			return err
		}
	}
	return nil
}

func (ci *WinterfellCi) Close() error {
	var msg []string
	if ci.Ipmi != nil {
		err := ci.Ipmi.Close()

		if err != nil {
			msg = append(msg, fmt.Sprintf("cannot close IPMI (err: %v)", err))
		}
	}
	if ci.Relay != nil {
		err := ci.Relay.PowerDown()
		if err != nil {
			if err != nil {
				msg = append(msg, fmt.Sprintf("cannot power down relay (err: %v)", err))
			}
		}
	}
	if ci.Em100 != nil {
		ci.Em100.Stop()
	}
	if ci.Serial != nil {
		err := ci.Serial.Close()

		if err != nil {
			msg = append(msg, fmt.Sprintf("cannot close serial (err: %v)", err))
		}
	}
	if len(msg) > 0 {
		return errors.New(strings.Join(msg, "; "))
	}
	return nil
}

func CloseWinterfellCiTest(t *testing.T, ci *WinterfellCi) {
	err := ci.Close()
	if err != nil {
		t.Fatal(err)
	}
}
