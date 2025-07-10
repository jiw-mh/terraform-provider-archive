// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

		isTemplate := false
		if opts.TemplateFileSuffix != "" {
			isTemplate = strings.HasSuffix(relname, opts.TemplateFileSuffix)
		}
		if isTemplate {
			relname = strings.TrimSuffix(relname, opts.TemplateFileSuffix)
		}

		archivePath := filepath.Join(basePath, relname)
		if opts.usedPaths[archivePath] {
			return fmt.Errorf("template path and non-template path found in same folder: %s", archivePath)
		}

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

		if isTemplate {
			tmpDir := ""
			tmpDir, path, info, err = renderTemplate(opts.Context, filepath.Base(relname), path, opts.TemplateVariables, info.Mode())
			if err != nil {
				return err
			}
			defer os.RemoveAll(tmpDir)
		}

		opts.usedPaths[archivePath] = true

		return process(path, archivePath, info)
	}
}

func renderTemplate(ctx context.Context, filename string, srcPath string, replacer func(string) (string, error), filemode os.FileMode) (string, string, os.FileInfo, error) {
	tflog.Info(ctx, "Rendering template", map[string]any{"srcPath": srcPath, "filename": filename})
	template, err := os.ReadFile(srcPath)
	if err != nil {
		return "", "", nil, fmt.Errorf("cannot read template in %s: %s", srcPath, err)
	}
	result, err := replacer(string(template))
	if err != nil {
		return "", "", nil, fmt.Errorf("error while processing template %s: %s", srcPath, err)
	}
	tmpDir, err := os.MkdirTemp("", "tf-provider-archive-*")
	if err != nil {
		return "", "", nil, fmt.Errorf("cannot create temp dir for %s: %s", srcPath, err)
	}
	tmpFile := filepath.Join(tmpDir, filename)
	err = os.WriteFile(tmpFile, []byte(result), filemode)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", "", nil, fmt.Errorf("cannot write temporary file to %s: %s", tmpFile, err)
	}
	info, err := os.Stat(tmpFile)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", "", nil, fmt.Errorf("cannot retreive stat from temporary file at %s: %s", tmpFile, err)
	}
	_, err = os.ReadFile(tmpFile)
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", "", nil, fmt.Errorf("cannot retreive temporary file: %s", err)
	}
	return tmpDir, tmpFile, info, nil
}

func assertArchiveDirHasFiles(indirname string, opts ArchiveDirOpts) error {
	isArchiveEmpty := true
	opts.usedPaths = make(map[string]bool)

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
	opts.usedPaths = make(map[string]bool)

	return filepath.Walk(indirname, createWalkFunc("", indirname, opts, &isArchiveEmpty, process))
}
