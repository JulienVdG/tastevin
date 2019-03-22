// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Csl open a console on a remote server via serial port
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/JulienVdG/tastevin/pkg/scriptreplay"
	"github.com/tarm/serial"
)

func main() {
	// open serial
	c := &serial.Config{Name: "/dev/ttyUSB0", Baud: 57600 /*, ReadTimeout: time.Nanosecond /*time.Second * 1.0 / 5760000*/}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := s.Close()

		if err != nil {
			fmt.Println(err)

		}
	}()
	// open output file
	fo, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}

	// open output file
	ft, err := os.Create("output.tim")
	if err != nil {
		panic(err)
	}

	sr := scriptreplay.NewWriter(fo, ft)
	defer func() {
		err := sr.Close()

		if err != nil {
			fmt.Println(err)

		}
	}()

	chanInToSerial := make(chan error)
	go func() {
		chanInToSerial <- copyStream(s, os.Stdin)
	}()
	chanSerialToOut := make(chan error)
	go func() {
		chanSerialToOut <- copyStream(io.MultiWriter(os.Stdout, sr), s)
	}()
	select {
	case err := <-chanSerialToOut:
		if err != nil {
			fmt.Println(err)
		}
		log.Println("connection closed")
	case err := <-chanInToSerial:
		if err != nil {
			fmt.Println(err)
		}
		log.Println("program terminated")
	}
}

func copyStream(dst io.Writer, src io.Reader) error {
	buf := make([]byte, 1024)
	var err error
	for {
		var n int
		n, err = src.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("Read error: %s\n", err)
			}
			break
		}
		_, err = dst.Write(buf[0:n])
		if err != nil {
			log.Fatalf("Write error: %s\n", err)
		}
	}
	return err
}
