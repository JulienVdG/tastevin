// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipmi

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// ConfigFromEnv returns the ipmi configuration from TASTEVIN_IPMI env var
// example env:
// TASTEVIN_IPMI='{"Hostname":"10.0.3.208","Username":"USERID","Password":"PASSW0RD","Interface":"lanplus","Path":"ipmitool","YafuFlashPath":"YafuFlash"}'
func ConfigFromEnv() (*Connection, error) {
	env := os.Getenv("TASTEVIN_IPMI")
	if env == "" {
		return nil, errors.New("TASTEVIN_IPMI not found in environment")
	}
	var c Connection
	err := json.Unmarshal([]byte(env), &c)
	if err != nil {
		return nil, fmt.Errorf("TASTEVIN_IPMI invalid %v", err)
	}
	if c.Interface == "" {
		c.Interface = "lanplus"
	}
	return &c, nil
}
