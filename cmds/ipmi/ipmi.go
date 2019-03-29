// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	ipmi "github.com/JulienVdG/tastevin/pkg/ipmi"
)

func main() {
	fmt.Println("vim-go")

	remote, err := ipmi.NewRemote(&ipmi.Connection{
		Hostname:  "192.168.142.187",
		Username:  "USERID",
		Password:  "PASSW0RD",
		Interface: "lanplus",
		Path:      "ipmitool",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	err = remote.Open()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	on, err := remote.PowerStatus()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Power status is %v", on)

	if on {
		remote.PowerDown()
	} else {
		remote.PowerUp()
	}

	err = remote.Close()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
