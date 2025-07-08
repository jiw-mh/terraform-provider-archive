// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type ZipArchiver struct {
	filepath       string
	outputFileMode string // Default value "" means unset
	filewriter     *os.File
	writer         *zip.Writer
}

func NewZipArchiver(filepath string) Archiver {
	return &ZipArchiver{
		filepath: filepath,
	}
}

func (a *ZipArchiver) ArchiveContent(content []byte, infilename string) error {
	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	f, err := a.writer.Create(filepath.ToSlash(infilename))
	if err != nil {
		return err
	}

	_, err = f.Write(content)
	return err
}

func (a *ZipArchiver) ArchiveFile(infilename string) error {
	fi, err := assertValidFile(infilename)
	if err != nil {
		return err
	}

	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	return a.addFileInfo(infilename, fi.Name(), fi)
}

func (a *ZipArchiver) ArchiveDir(indirname string, opts ArchiveDirOpts) error {
	err := assertArchiveDirHasFiles(indirname, opts)
	if err != nil {
		return err
	}

	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	return walkDir(indirname, opts, a.addFileInfo)
}

func (a *ZipArchiver) addFileInfo(path string, archivePath string, info os.FileInfo) error {
	fh, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("error creating file header: %s", err)
	}
	fh.Name = filepath.ToSlash(archivePath)
	fh.Method = zip.Deflate
	// fh.Modified alone isn't enough when using a zero value
	//nolint:staticcheck
	fh.SetModTime(time.Time{})

	if a.outputFileMode != "" {
		filemode, err := strconv.ParseUint(a.outputFileMode, 0, 32)
		if err != nil {
			return fmt.Errorf("error parsing output_file_mode value: %s", a.outputFileMode)
		}
		fh.SetMode(os.FileMode(filemode))
	}

	f, err := a.writer.CreateHeader(fh)
	if err != nil {
		return fmt.Errorf("error creating file inside archive: %s", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file for archival: %s", err)
	}
	_, err = f.Write(content)
	return err
}

func (a *ZipArchiver) ArchiveMultiple(content map[string][]byte) error {
	if err := a.open(); err != nil {
		return err
	}
	defer a.close()

	// Ensure files are processed in the same order so hashes don't change
	keys := make([]string, len(content))
	i := 0
	for k := range content {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, filename := range keys {
		f, err := a.writer.Create(filepath.ToSlash(filename))
		if err != nil {
			return err
		}
		_, err = f.Write(content[filename])
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ZipArchiver) SetOutputFileMode(outputFileMode string) {
	a.outputFileMode = outputFileMode
}

func (a *ZipArchiver) open() error {
	f, err := os.Create(a.filepath)
	if err != nil {
		return err
	}
	a.filewriter = f
	a.writer = zip.NewWriter(f)
	return nil
}

func (a *ZipArchiver) close() {
	if a.writer != nil {
		a.writer.Close()
		a.writer = nil
	}
	if a.filewriter != nil {
		a.filewriter.Close()
		a.filewriter = nil
	}
}
