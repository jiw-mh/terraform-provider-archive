// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

func checkMatch(fileName string, excludes []string) (value bool, err error) {
	for _, exclude := range excludes {
		if exclude == "" {
			continue
		}

		match, err := doublestar.PathMatch(exclude, fileName)
		if err != nil {
			return false, err
		}

		if match {
			return true, nil
		}
	}
	return false, nil
}

func createWalkFunc(basePath string, indirname string, opts ArchiveDirOpts, isArchiveEmpty *bool, process func(path string, archivePath string, info os.FileInfo) error) func(path string, info os.FileInfo, err error) error {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error encountered during file walk: %s", err)
		}

		relname, err := filepath.Rel(indirname, path)
		if err != nil {
			return fmt.Errorf("error relativizing file for archival: %s", err)
		}

		archivePath := filepath.Join(basePath, relname)

		isMatch, err := checkMatch(archivePath, opts.Excludes)
		if err != nil {
			return fmt.Errorf("error checking excludes matches: %w", err)
		}

		if info.IsDir() {
			if isMatch {
				return filepath.SkipDir
			}
			return nil
		}

		if isMatch {
			return nil
		}

		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				return err
			}

			realInfo, err := os.Stat(realPath)
			if err != nil {
				return err
			}

			if realInfo.IsDir() {
				if !opts.ExcludeSymlinkDirectories {
					return filepath.Walk(realPath, createWalkFunc(archivePath, realPath, opts, isArchiveEmpty, process))
				} else {
					return filepath.SkipDir
				}
			}

			info = realInfo
		}

		*isArchiveEmpty = false

		if process == nil {
			return nil
		}

		return process(path, archivePath, info)
	}
}

func assertArchiveDirHasFiles(indirname string, opts ArchiveDirOpts) error {
	isArchiveEmpty := true

	err := filepath.Walk(indirname, createWalkFunc("", indirname, opts, &isArchiveEmpty, nil))

	if err != nil {
		return err
	}

	// Return an error if an empty archive would be generated.
	if isArchiveEmpty {
		return fmt.Errorf("archive has not been created as it would be empty")
	}
	return nil
}

func walkDir(indirname string, opts ArchiveDirOpts, process func(path string, archivePath string, info os.FileInfo) error) error {
	// Needed for implementation
	isArchiveEmpty := true

	return filepath.Walk(indirname, createWalkFunc("", indirname, opts, &isArchiveEmpty, process))
}
