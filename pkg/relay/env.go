// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package relay

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type relayConf struct {
	PowerOnCmd  []string
	PowerOffCmd []string
}

// NewRelayFromEnv returns a Relay from TASTEVIN_RELAY env var
// example env:
// TASTEVIN_RELAY='{"PowerOnCmd":["sudo","/usr/local/bin/relay-poweron"],"PowerOffCmd":["sudo","/usr/local/bin/relay-poweroff"]}'
func NewRelayFromEnv() (*Relay, error) {
	env := os.Getenv("TASTEVIN_RELAY")
	if env == "" {
		return nil, errors.New("TASTEVIN_RELAY not found in environment")
	}
	var c relayConf
	err := json.Unmarshal([]byte(env), &c)
	if err != nil {
		return nil, fmt.Errorf("TASTEVIN_RELAY invalid %v", err)
	}
	if len(c.PowerOnCmd) == 0 {
		return nil, fmt.Errorf("TASTEVIN_RELAY missing PowerOnCmd")
	}
	if len(c.PowerOffCmd) == 0 {
		return nil, fmt.Errorf("TASTEVIN_RELAY missing PowerOffCmd")
	}
	r := Relay{
		powerOnArgs:  c.PowerOnCmd,
		powerOffArgs: c.PowerOffCmd,
	}
	return &r, nil
}
