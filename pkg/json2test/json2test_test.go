// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package json2test_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/JulienVdG/tastevin/pkg/json2test"
	"github.com/JulienVdG/tastevin/pkg/xio/iotest"
)

func Assert(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

type testCase struct {
	iotest.CallCounter
	e      json2test.TestEvent
	result error
}

func (c *testCase) Handle(e json2test.TestEvent) {
	c.CallCount++
	if !reflect.DeepEqual(e, c.e) {
		c.result = fmt.Errorf("TestEvent differ got %v expected %v", e, c.e)
	}
	c.result = nil
}

var logCount iotest.CallCounter

func logCounter(string, ...interface{}) {
	logCount.CallCount++
}

func TestConverter(t *testing.T) {
	json2test.LogWarn = logCounter
	c := &testCase{}
	w := json2test.NewConverter(c)
	t.Run("Empty", func(t *testing.T) {
		i, err := w.Write([]byte{})
		Assert(t, err)
		if i != 0 {
			t.Errorf("expected 0 got %v", i)
		}
		Assert(t, c.NotCalled())
		Assert(t, logCount.NotCalled())
	})
	t.Run("InvalidJSON", func(t *testing.T) {
		i, err := w.Write([]byte("}"))
		Assert(t, err)
		if i != 1 {
			t.Errorf("expected 1 got %v", i)
		}
		Assert(t, logCount.CalledOnce())
		Assert(t, c.NotCalled())
	})
	time0 := time.Unix(0, 0)
	elapsed0 := float64(0.0)
	var tests = []struct {
		name string
		e    json2test.TestEvent
	}{
		{"Nil", json2test.TestEvent{}},
		{"case1", json2test.TestEvent{&time0, "action", "package", "test", &elapsed0, "Output"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c.e = test.e
			msg, err := json.Marshal(&test.e)
			Assert(t, err)
			i, err := w.Write(msg)
			Assert(t, err)
			if i != len(msg) {
				t.Errorf("expected %d got %d", len(msg), i)
			}
			Assert(t, c.result)
			Assert(t, c.CalledOnce())
			Assert(t, logCount.NotCalled())
		})

	}
}

func TestConverterVerbose(t *testing.T) {
	json2test.LogWarn = logCounter

	files, err := filepath.Glob("testdata/*.test")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		name := strings.TrimSuffix(filepath.Base(file), ".test")
		t.Run(name, func(t *testing.T) {
			orig, err := ioutil.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}

			var buf bytes.Buffer
			in := orig

			c := json2test.NewConverter(json2test.NewVerboseHandler(&buf))
			cmd := exec.Command("go", "tool", "test2json")
			cmd.Stdin = bytes.NewBuffer(in)
			cmd.Stdout = c
			cmd.Stderr = c
			err = cmd.Run()
			Assert(t, err)
			Assert(t, logCount.NotCalled())

			res := buf.Bytes()
			if bytes.Compare(orig, res) != 0 {
				t.Errorf("Content differ got:\n%q\nexpect:\n%q\n", res, orig)
			}

		})
	}

}
