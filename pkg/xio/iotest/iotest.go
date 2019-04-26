// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package iotest provide helpers to test io interfaces
package iotest

import (
	"bytes"
	"fmt"
)

// CallCounter is used to track the number of function call.
type CallCounter struct {
	Name      string
	CallCount int
}

// ExpectWriter provides the io.Writer interface, count the Write calls and
// check the write content.
type ExpectWriter struct {
	CallCounter
	Content []byte
	Result  error
}

// ExpectCloser provides the io.Closer interface, count the Close calls.
type ExpectCloser struct {
	CallCounter
}

// ExpectWriteCloser provides both the io.Writer and the io.Closer interfaces,
// count the Write and Close calls and check the write content.
type ExpectWriteCloser struct {
	ExpectWriter
	ExpectCloser
}

// Write provide io.Write interface on ExpectWriter _pointer_.
func (ew *ExpectWriter) Write(p []byte) (n int, err error) {
	ew.CallCount++
	cmp := bytes.Compare(ew.Content, p)
	if cmp != 0 {
		ew.Result = fmt.Errorf("mismatched %s write, wanted '%v', got '%v'", ew.Name, string(ew.Content), string(p))
	}
	return len(p), nil
}

// Close provide io.Close interface on ExpectCloser _pointer_.
func (ec *ExpectCloser) Close() error {
	ec.CallCount++
	return nil
}

// NotCalled return nil if not called
func (cc *CallCounter) NotCalled() (err error) {
	if cc.CallCount != 0 {
		err = fmt.Errorf("%s called while not expected, got %v", cc.Name, cc.CallCount)
	}
	cc.CallCount = 0
	return
}

// CalledOnce return nil if called exactly once
func (cc *CallCounter) CalledOnce() (err error) {
	if cc.CallCount == 0 {
		err = fmt.Errorf("%s not called, expected once", cc.Name)
	} else if cc.CallCount != 1 {
		err = fmt.Errorf("%s called more than once, got %v", cc.Name, cc.CallCount)
	}
	cc.CallCount = 0
	return
}
