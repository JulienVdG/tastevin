// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xio_test

import (
	"fmt"
	"io"
	"testing"

	"github.com/JulienVdG/tastevin/pkg/xio"
	"github.com/JulienVdG/tastevin/pkg/xio/iotest"
)

func Assert(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestCloser(t *testing.T) {

	var ecs [2]iotest.ExpectCloser
	cs := make([]io.Closer, len(ecs))
	for i, _ := range ecs {
		ecs[i].Name = fmt.Sprintf("closer %d", i)
		cs[i] = &ecs[i]
	}

	mc := xio.MultiCloser(cs...)

	for i, _ := range ecs {
		Assert(t, ecs[i].NotCalled())
	}

	mc.Close()

	for i, _ := range ecs {
		Assert(t, ecs[i].CalledOnce())
	}
}

func TestWriteCloser(t *testing.T) {

	var ewcs [2]iotest.ExpectWriteCloser
	wcs := make([]io.WriteCloser, len(ewcs))
	for i, _ := range ewcs {
		ewcs[i].ExpectWriter.Name = fmt.Sprintf("writer %d", i)
		ewcs[i].ExpectCloser.Name = fmt.Sprintf("closer %d", i)
		wcs[i] = &ewcs[i]
	}

	mwc := xio.MultiWriteCloser(wcs...)

	for i, _ := range ewcs {
		Assert(t, ewcs[i].ExpectWriter.NotCalled())
		Assert(t, ewcs[i].ExpectCloser.NotCalled())
	}

	var tests = []string{
		"Send some char\n",
		"more lines\n...\neven more",
	}
	for _, test := range tests {
		for i, _ := range ewcs {
			ewcs[i].Content = []byte(test)
		}

		mwc.Write([]byte(test))

		for i, _ := range ewcs {
			Assert(t, ewcs[i].ExpectWriter.Result)
			Assert(t, ewcs[i].ExpectWriter.CalledOnce())
			Assert(t, ewcs[i].ExpectCloser.NotCalled())
		}
	}

	mwc.Close()

	for i, _ := range ewcs {
		Assert(t, ewcs[i].ExpectWriter.NotCalled())
		Assert(t, ewcs[i].ExpectCloser.CalledOnce())
	}
}

func TestCloser_merged(t *testing.T) {

	var ecs [3]iotest.ExpectCloser
	cs := make([]io.Closer, len(ecs))
	for i, _ := range ecs {
		ecs[i].Name = fmt.Sprintf("closer %d", i)
		cs[i] = &ecs[i]
	}

	mc1 := xio.MultiCloser(cs[1:]...)
	mc := xio.MultiCloser(cs[0], mc1)

	for i, _ := range ecs {
		Assert(t, ecs[i].NotCalled())
	}

	mc.Close()

	for i, _ := range ecs {
		Assert(t, ecs[i].CalledOnce())
	}
}
