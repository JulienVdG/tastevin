// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// winterfell-ci controls a winterfell node with relay and em100.

// TODO handle errors

package main

import (
	"flag"
	"fmt"

	"github.com/JulienVdG/tastevin/pkg/em100"
	"github.com/JulienVdG/tastevin/pkg/relay"
)

var (
	dumpFile  = flag.String("dump", "", "dump the emulator content to file (will also stop the node)")
	loadFile  = flag.String("load", "", "load a flash firmware into the emulator (will restart the node)")
	startFlag = flag.Bool("start", false, "start the node (if started it will restart it)")
	stopFlag  = flag.Bool("stop", false, "stop the node")
)

func main() {
	flag.Parse()

	// Emulator commands
	doDump := *dumpFile != ""
	doLoad := *loadFile != ""
	willStopEm := doDump || doLoad
	// Relay + emulator commands
	doStop := *stopFlag || willStopEm
	doStart := *startFlag || doLoad

	// Create controlling objects
	em, err := em100.NewEm100FromEnv()
	if err != nil {
		em = em100.NewEm100("", "")
	}
	r, err := relay.NewRelayFromEnv()
	if err != nil {
		r = relay.NewRelay()
	}

	// Process command in order
	if doStop {
		fmt.Println("Stopping node")
		r.PowerDown()
		if !willStopEm {
			em.Stop()
		}
	}
	if doDump {
		fmt.Printf("Dumping flash emulator content to '%s'\n", *dumpFile)
		em.Dump(*dumpFile)
	}
	if doLoad {
		fmt.Printf("Loading flash emulator content from '%s'\n", *loadFile)
		em.Load(*loadFile)
	}
	if doStart {
		if !doLoad {
			em.Start()
		}
		fmt.Println("Starting node")
		r.PowerUp()
	}
}
