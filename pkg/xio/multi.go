// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xio

import (
	"errors"
	"io"
	"strings"
)

type multiCloser struct {
	closers []io.Closer
}

// Close closes all closers of MultiCloser
func (mc *multiCloser) Close() error {
	var msg []string
	for _, c := range mc.closers {
		if err := c.Close(); err != nil {
			msg = append(msg, err.Error())
		}
	}
	if len(msg) > 0 {
		return errors.New(strings.Join(msg, "; "))
	}
	return nil
}

// MultiCloser creates a closer that will close all the provided closers.
func MultiCloser(closers ...io.Closer) io.Closer {
	allClosers := make([]io.Closer, 0, len(closers))
	for _, c := range closers {
		if mc, ok := c.(*multiCloser); ok {
			allClosers = append(allClosers, mc.closers...)
		} else {
			allClosers = append(allClosers, c)
		}
	}
	return &multiCloser{allClosers}
}

type multiWriteCloser struct {
	io.Writer
	io.Closer
}

// MultiWriteCloser create a WriteCloser that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command and close them all on
// Close.
func MultiWriteCloser(writeClosers ...io.WriteCloser) io.WriteCloser {
	allWriters := make([]io.Writer, len(writeClosers))
	allClosers := make([]io.Closer, len(writeClosers))
	for i, wc := range writeClosers {
		allWriters[i] = wc
		allClosers[i] = wc
	}
	w := io.MultiWriter(allWriters...)
	c := MultiCloser(allClosers...)

	return &multiWriteCloser{w, c}
}
