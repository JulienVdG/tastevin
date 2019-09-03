// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/JulienVdG/tastevin/pkg/gotestweb"
)

// Gen generate an output directory suitable for serving GoTestWeb and the logs on an external web server.
// outdir         : absolute output directory
// logdir         : absolute log directory,
//                  if outside of outdir, it will be copied in
// jsonPath       : path to json file inside logdir
// indexPath      : index name inside outdir
// externalAppURL : path to the server URL of GoTestWeb,
//                  if empty the GoTestWeb app will be copied to outdir
// useCDN         : whether to use public CDN URL for bootstrap and jquery"
func Gen(outdir, logdir, jsonPath, indexPath, externalAppURL string, useCDN bool) error {
	var slug, jsonfile string
	if !strings.HasPrefix(logdir+"/", outdir+"/") {
		slug = filepath.Base(logdir)
		jsonfile = slug + "/" + jsonPath
		// copy logs if logdir is not inside output dir
		outlogdir := filepath.Join(outdir, slug)
		e := filepath.Walk(logdir, func(path string, f os.FileInfo, err error) error {
			if err != nil { // Don't try to fix walk issues
				return err
			}
			relpath, err := filepath.Rel(logdir, path)
			if err != nil {
				return err
			}
			target := filepath.Join(outlogdir, relpath)
			return copyItemTo(path, target, f, nil)
		})
		if e != nil {
			return e
		}
	} else {
		// logdir is inside output dir, rebuild path
		var err error
		slug, err = filepath.Rel(outdir+"/", logdir+"/")
		if err != nil {
			return err
		}

		if slug == "." && jsonPath == filepath.Base(jsonPath) {
			slug = ""
			jsonfile = jsonPath
		} else {
			jsonfile = slug + "/" + jsonPath
		}
	}

	// use template to generate index
	indexfilename := filepath.Join(outdir, indexPath)
	f, err := os.Create(indexfilename)
	if err != nil {
		return fmt.Errorf("error creating file '%s': %v", indexPath, err)
	}
	defer f.Close()
	data := gotestweb.IndexData{
		File:         jsonfile,
		Asciicast:    slug,
		Scriptreplay: slug,
		AppPrefix:    externalAppURL,
		UseCDN:       useCDN,
	}
	gotestweb.WriteIndex(f, data)

	if externalAppURL == "" {
		box, err := gotestweb.RiceBox()
		if err != nil {
			return err
		}
		e := box.Walk("", func(path string, f os.FileInfo, err error) error {
			if err != nil { // Don't try to fix walk issues
				return err
			}
			// skip index, we generated it anyway
			if path == "index.html" {
				return nil
			}
			// skip box root
			if f.Name() == "http-files" {
				return nil
			}
			if useCDN {
				// skip cdn components
				switch path {
				case "vendor/bootstrap", "vendor/jquery":
					return filepath.SkipDir
				}
			}
			target := filepath.Join(outdir, path)
			return copyItemTo(path, target, f, box)
		})
		if e != nil {
			return e
		}
	}
	return nil
}

func copyItemTo(src, dst string, srcfi os.FileInfo, box *rice.Box) error {
	m := srcfi.Mode()
	//fmt.Println(src, dst, m.IsDir(), m.Perm())
	if m.IsDir() {
		return os.MkdirAll(dst, 0775)
	}
	if !m.IsRegular() {
		return fmt.Errorf("unsupported file mode %s", m)
	}
	return copyRegularFile(src, dst, srcfi, box)
}

func copyRegularFile(src, dst string, srcfi os.FileInfo, box *rice.Box) error {
	var srcf io.Reader
	if box == nil {
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()
		srcf = f
	} else {
		f, err := box.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()
		srcf = f
	}

	dstf, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, srcfi.Mode().Perm())
	if err != nil {
		return err
	}
	defer dstf.Close()

	_, err = io.Copy(dstf, srcf)
	return err
}
