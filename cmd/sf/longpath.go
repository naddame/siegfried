// +build !windows

// Copyright 2015 Richard Lehane. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/richardlehane/siegfried"
)

func retryOpen(path string, err error) (*os.File, error) {
	return nil, err
}

func retryStat(path string, err error) (os.FileInfo, error) {
	return nil, err
}

func identify(ctxts chan *context, wg *sync.WaitGroup, s *siegfried.Siegfried, w writer, root, orig string, h hash.Hash, z bool, norecurse bool) error {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if *throttlef > 0 {
			<-throttle.C
		}
		if err != nil {
			return fmt.Errorf("walking %s; got %v", path, err)
		}
		if info.IsDir() {
			if norecurse && path != root {
				return filepath.SkipDir
			}
			if *droido {
				dctx := newContext(s, w, wg, nil, false, path, "", info.ModTime().Format(time.RFC3339), -1)
				dctx.res <- results{nil, nil, nil}
				wg.Add(1)
				ctxts <- dctx
			}
			return nil
		}
		identifyFile(newContext(s, w, wg, h, z, path, "", info.ModTime().Format(time.RFC3339), info.Size()), ctxts)
		return nil
	}
	return filepath.Walk(root, walkFunc)
}
