// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	goipmi "github.com/vmware/goipmi"
)

func TestRemote(t *testing.T) {
	s := goipmi.NewSimulator(net.UDPAddr{Port: 0})
	err := s.Run()
	assert.NoError(t, err)

	c := Connection(*s.NewConnection())
	remote, err := NewRemote(&c)
	assert.NoError(t, err)

	err = remote.Open()
	assert.NoError(t, err)

	on, err := remote.PowerStatus()
	assert.NoError(t, err)
	// t.Logf("Power status is %v", on)
	assert.True(t, on)

	err = remote.Close()
	assert.NoError(t, err)
	s.Stop()
}
