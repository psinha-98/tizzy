package tizzy

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"archive/zip"
	"errors"
	"log"
	"runtime"
	"time"
)

var (
	timezones = map[string]time.Location{}
)

func init() {
	zoneinfo := runtime.GOROOT() + "/lib/time/zoneinfo.zip"

	r, err := zip.OpenReader(zoneinfo)
	if err != nil {
		log.Fatalf("Could not open zoneinfo.zip for reading: path '%s' error: '%v'", zoneinfo, err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f == nil {
			log.Fatalf("zoneinfo.zip contained nil file pointer")
		}

		if f.FileHeader.FileInfo().IsDir() {
			continue
		}

		if _, ok := timezones[f.Name]; ok {
			continue
		}

		tz, err := time.LoadLocation(f.Name)
		if err != nil {
			log.Fatalf("Could not load location '%s': '%v'", f.Name, err)
		}

		timezones[f.Name] = *tz
	}
}

// LoadLocation returns the Location with the given name.
//
// If the name is "" or "UTC", LoadLocation returns UTC.
// If the name is "Local", LoadLocation returns Local.
//
// Otherwise, the name is taken to be a location name corresponding to a file
// in the IANA Time Zone database, such as "America/New_York".
func LoadLocation(name string) (*time.Location, error) {
	if name == "" || name == "UTC" {
		return time.UTC, nil
	}

	if name == "Local" {
		return time.Local, nil
	}

	if containsDotDot(name) || name[0] == '/' || name[0] == '\\' {
		// No valid IANA Time Zone name contains a single dot,
		// much less dot dot. Likewise, none begin with a slash.
		return nil, errors.New("time: invalid location name")
	}

	zoneData, ok := timezones[name]
	if !ok {
		return nil, errors.New("unknown time zone " + name)
	}

	return &zoneData, nil
}

// containsDotDot reports whether s contains "..".
func containsDotDot(s string) bool {
	if len(s) < 2 {
		return false
	}
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '.' && s[i+1] == '.' {
			return true
		}
	}
	return false
}
