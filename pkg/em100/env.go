// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package em100

import (
	"errors"
	"os"
	"strings"
)

// NewEm100FromEnv created a new Em100 using TASTEVIN_EM100 env var
// example env:
// TASTEVIN_EM100='sudo /opt/em100/em100 --set W25Q128BV'
func NewEm100FromEnv() (*Em100, error) {
	env := os.Getenv("TASTEVIN_EM100")
	if env == "" {
		return nil, errors.New("TASTEVIN_EM100 not found in environment")
	}
	var args []string
	args = append(args, strings.Fields(env)...)
	return &Em100{args: args}, nil
}
