// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package main

import (
	"github.com/alexandr-andreyev/soup-rk7-events/internal/app"
	"github.com/pkg/errors"
)

// This is the name you will use for the NET START command
const svcName = "soup-rk7-events" //TODO Имя службы брать из настроек

// This is the name that will appear in the Services control panel
const svcNameLong = "Soup Events for R_Keeper7" //TODO Имя службы брать из настроек

// This is assigned the full SHA1 hash from GIT
var sha1ver string

func svcLauncher() error {

	err := app.Run(elog, svcName, sha1ver)
	if err != nil {
		return errors.Wrap(err, "app.run")
	}

	return nil
}
