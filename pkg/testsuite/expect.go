// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testsuite

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/JulienVdG/tastevin/pkg/asciicast"
	"github.com/JulienVdG/tastevin/pkg/scriptreplay"
	"github.com/JulienVdG/tastevin/pkg/xio"
	exp "github.com/google/goexpect"
)

func openScriptReplay(prefix, baseName string) (io.WriteCloser, error) {
	if len(prefix) == 0 {
		return nil, nil
	}
	dir := filepath.Dir(prefix)
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return nil, err
	}
	f := prefix + baseName
	return scriptreplay.NewFileWriter(f+".log", f+".tim")
}

func openAsciicast(prefix, baseName string) (io.WriteCloser, error) {
	if len(prefix) == 0 {
		return nil, nil
	}
	dir := filepath.Dir(prefix)
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return nil, err
	}
	return asciicast.NewFileWriter(prefix + baseName + ".cast")
}

// ExpectOptions return options based on configuration from env
// returned error cn be considered a warning and logged only
func ExpectOptions(screenLogBaseName string) ([]exp.Option, error) {
	var msg []string
	o := []exp.Option{exp.PartialMatch(true)}
	c, err := getConf()
	if err != nil {
		msg = append(msg, fmt.Sprintf("%v, using default config", err))
	}
	if c.ExpectDebugCheck {
		o = append(o, exp.DebugCheck(nil))
	}
	if c.ExpectVerbose {
		o = append(o, exp.Verbose(true))
	}
	if len(screenLogBaseName) == 0 {
		screenLogBaseName = callerName(2, c.LongName)
	}
	swc, err := openScriptReplay(c.ScriptReplayPrefix, screenLogBaseName)
	if err != nil {
		msg = append(msg, fmt.Sprintf("skipping ScriptReplay output (err:%v)", err))
	}
	awc, err := openAsciicast(c.AsciicastPrefix, screenLogBaseName)
	if err != nil {
		msg = append(msg, fmt.Sprintf("skipping asciicast output (err:%v)", err))
	}
	if swc != nil && awc != nil {
		wc := xio.MultiWriteCloser(swc, awc)
		o = append(o, exp.Tee(wc))
	} else if swc != nil {
		o = append(o, exp.Tee(swc))
	} else if awc != nil {
		o = append(o, exp.Tee(awc))
	}

	if len(msg) > 0 {
		return o, errors.New(strings.Join(msg, "; "))
	}
	return o, nil
}

func DescribeBatcherErr(batch []exp.Batcher, res []exp.BatchRes, err error) error {
	last := res[len(res)-1]
	i := last.Idx
	return fmt.Errorf("Batch[%d]:%+v err: %v", i, batch[i], err)
}
