// Copyright 2019 Splitted-Desktop Systems. All rights reserved
// Copyright 2019 Julien Viard de Galbert
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gotestweb

import (
	"net/http"

	rice "github.com/GeertJohan/go.rice"
)

// Handle adds a http.Handle for / serving the gotest-web js application
func Handle() error {
	box, err := rice.FindBox("http-files")
	if err != nil {
		return err
	}
	http.Handle("/", http.FileServer(box.HTTPBox()))
	return nil
}
