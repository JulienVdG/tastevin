// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testsuite

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// TASTEVIN_CONF='{"ScriptReplayPrefix":"path/to/logprefix-","AsciicastPrefix":"path/to/logprefix-","ExpectDebugCheck":true,"ExpectVerbose":true,"LongName":true}'
type conf struct {
	ScriptReplayPrefix string
	AsciicastPrefix    string
	ExpectDebugCheck   bool
	ExpectVerbose      bool
	LongName           bool
	loaded             bool
}

var _conf conf

func getConf() (*conf, error) {
	if _conf.loaded {
		return &_conf, nil
	}
	_conf.loaded = true
	env := os.Getenv("TASTEVIN_CONF")
	if env == "" {
		return &_conf, errors.New("TASTEVIN_CONF not found in environment")
	}
	var c conf
	err := json.Unmarshal([]byte(env), &c)
	if err != nil {
		return &_conf, fmt.Errorf("error parsing TASTEVIN_CONF: %v", err)
	}
	_conf = c

	return &_conf, nil
}
