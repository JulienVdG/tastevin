// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ipmi implement the XXX interface over IPMI
package ipmi

import (
	"fmt"

	goipmi "github.com/vmware/goipmi"
)

// Connection properties for a Client
type Connection struct {
	goipmi.Connection
	YafuFlashPath string `json:",omitempty"`
}

// Remote represent an IPMI Remote
type Remote struct {
	*goipmi.Client
	YafuFlashPath string
}

// NewRemote creates a new IPMI Remote
func NewRemote(c *Connection) (*Remote, error) {
	cc := c.Connection
	client, err := goipmi.NewClient(&cc)
	return &Remote{Client: client, YafuFlashPath: c.YafuFlashPath}, err
}

// Open a new IPMI session
func (r *Remote) Open() error {
	return r.Client.Open()
}

// Close the IPMI session
func (r *Remote) Close() error {
	return r.Client.Close()
}

// PowerStatus returns the current power status
func (r *Remote) PowerStatus() (bool, error) {
	req := &goipmi.Request{
		NetworkFunction: goipmi.NetworkFunctionChassis,
		Command:         goipmi.CommandChassisStatus,
		Data:            &goipmi.DeviceIDRequest{},
	}
	csr := &goipmi.ChassisStatusResponse{}
	err := r.Client.Send(req, csr)
	if err != nil {
		return false, fmt.Errorf("error requesting power status: %v", err)
	}
	return (csr.PowerState & goipmi.SystemPower) != 0, nil
}

// PowerUp sends the request to power up
func (r *Remote) PowerUp() error {
	on, err := r.PowerStatus()
	if err != nil {
		return err
	}
	if on {
		return nil
	}
	err = r.Client.Control(goipmi.ControlPowerUp)
	if err != nil {
		return fmt.Errorf("error powering up: %v", err)
	}
	return nil

}

// PowerDown sends the request to power down
func (r *Remote) PowerDown() error {
	on, err := r.PowerStatus()
	if err != nil {
		return err
	}
	if !on {
		return nil
	}
	err = r.Client.Control(goipmi.ControlPowerDown)
	if err != nil {
		return fmt.Errorf("error powering down: %v", err)
	}
	return nil
}
