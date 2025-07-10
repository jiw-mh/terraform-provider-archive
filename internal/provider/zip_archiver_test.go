// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestZipArchiver_Content(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-content.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveContent([]byte("This is some content"), "content.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"content.txt": []byte("This is some content"),
	})
}

func TestZipArchiver_File(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-file.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveFile("./test-fixtures/test-dir/test-file.txt"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"test-file.txt": []byte("This is test content"),
	})
}

//nolint:usetesting
func TestZipArchiver_FileMode(t *testing.T) {
	file, err := os.CreateTemp("", "archive-file-mode-test.zip")
	if err != nil {
		t.Fatal(err)
	}

	var (
		zipFilePath = file.Name()
		toZipPath   = filepath.FromSlash("./test-fixtures/test-dir/test-file.txt")
	)

	stringArray := [5]string{"0444", "0644", "0666", "0744", "0777"}
	for _, element := range stringArray {
		archiver := NewZipArchiver(zipFilePath)
		archiver.SetOutputFileMode(element)
		if err := archiver.ArchiveFile(toZipPath); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		ensureFileMode(t, zipFilePath, element)
	}
}

func TestZipArchiver_FileModified(t *testing.T) {
	var (
		zipFilePath = filepath.Join(t.TempDir(), "archive-file-modified.zip")
		toZipPath   = filepath.FromSlash("./test-fixtures/test-dir/test-file.txt")
	)

	var zipFunc = func() {
		archiver := NewZipArchiver(zipFilePath)
		if err := archiver.ArchiveFile(toZipPath); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	}

	zipFunc()

	expectedContents, err := os.ReadFile(zipFilePath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	//touch file modified, in the future just in case of weird race issues
	newTime := time.Now().Add(1 * time.Hour)
	if err := os.Chtimes(toZipPath, newTime, newTime); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	zipFunc()

	actualContents, err := os.ReadFile(zipFilePath)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if !bytes.Equal(expectedContents, actualContents) {
		t.Fatalf("zip contents do not match, potentially a modified time issue")
	}
}

func TestZipArchiver_Dir(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir/test-dir1", ArchiveDirOpts{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	})
}

func TestZipArchiver_Dir_Exclude(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-exclude.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir/test-dir1", ArchiveDirOpts{
		Excludes: []string{"file2.txt"},
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file3.txt": []byte("This is file 3"),
	})
}

func TestZipArchiver_Dir_Exclude_With_Directory(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-exclude-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir", ArchiveDirOpts{
		Excludes: []string{"test-dir1", "test-dir2/file2.txt"},
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"test-dir2/file1.txt": []byte("This is file 1"),
		"test-dir2/file3.txt": []byte("This is file 3"),
		"test-file.txt":       []byte("This is test content"),
	})
}

func TestZipArchiver_Multiple(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-content.zip")

	content := map[string][]byte{
		"file1.txt": []byte("This is file 1"),
		"file2.txt": []byte("This is file 2"),
		"file3.txt": []byte("This is file 3"),
	}

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveMultiple(content); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, content)
}

func TestZipArchiver_Dir_With_Symlink_File(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-file.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures/test-dir-with-symlink-file", ArchiveDirOpts{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"test-file.txt":    []byte("This is test content"),
		"test-symlink.txt": []byte("This is test content"),
	})
}

func TestZipArchiver_Dir_DoNotExcludeSymlinkDirectories(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures", ArchiveDirOpts{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ensureContents(t, zipFilePath, allFixturesInclSymlinks())
}

func TestZipArchiver_Dir_ExcludeSymlinkDirectories(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	err := archiver.ArchiveDir("./test-fixtures", ArchiveDirOpts{
		ExcludeSymlinkDirectories: true,
	})

	if err != nil {
		t.Errorf("expected no error: %s", err)
	}
}

func TestZipArchiver_Dir_Exclude_DoNotExcludeSymlinkDirectories(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures", ArchiveDirOpts{
		Excludes: []string{
			"test-symlink-dir/file1.txt",
			"test-symlink-dir-with-symlink-file/test-symlink.txt",
		},
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	wants := allFixturesInclSymlinks()
	delete(wants, "test-symlink-dir/file1.txt")
	delete(wants, "test-symlink-dir-with-symlink-file/test-symlink.txt")
	ensureContents(t, zipFilePath, wants)
}

func TestZipArchiver_Dir_Exclude_Glob_DoNotExcludeSymlinkDirectories(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	if err := archiver.ArchiveDir("./test-fixtures", ArchiveDirOpts{
		Excludes: []string{
			"**/file1.txt",
			"**/file2.txt",
			"test-dir-with-symlink-dir/test-symlink-dir",
			"test-symlink-dir-with-symlink-file/test-symlink.txt",
		},
	}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	wants := allFixturesInclSymlinks()
	delete(wants, "test-symlink-dir-with-symlink-file/test-symlink.txt")
	for key := range wants {
		if strings.HasSuffix(key, "file1.txt") ||
			strings.HasSuffix(key, "file2.txt") ||
			strings.HasPrefix(key, "test-dir-with-symlink-dir/test-symlink-dir") {
			delete(wants, key)
		}
	}
	ensureContents(t, zipFilePath, wants)
}

func TestZipArchiver_Dir_Exclude_ExcludeSymlinkDirectories(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	err := archiver.ArchiveDir("./test-fixtures", ArchiveDirOpts{
		Excludes: []string{
			"test-dir/test-dir1/file1.txt",
			"test-symlink-dir-with-symlink-file/test-symlink.txt",
		},
		ExcludeSymlinkDirectories: true,
	})

	if err != nil {
		t.Errorf("expected no error: %s", err)
	}

	wants := allFixtures()
	delete(wants, "test-dir/test-dir1/file1.txt")
	delete(wants, "test-symlink-dir-with-symlink-file/test-symlink.txt")
	ensureContents(t, zipFilePath, wants)
}

func TestZipArchiver_Dir_Exclude_Glob_ExcludeSymlinkDirectories(t *testing.T) {
	zipFilePath := filepath.Join(t.TempDir(), "archive-dir-with-symlink-dir.zip")

	archiver := NewZipArchiver(zipFilePath)
	err := archiver.ArchiveDir("./test-fixtures", ArchiveDirOpts{
		Excludes: []string{
			"test-dir/test-dir1/file1.txt",
			"**/file[2-3].txt",
			"test-dir-with-symlink-file",
		},
		ExcludeSymlinkDirectories: true,
	})

	if err != nil {
		t.Errorf("expected no error: %s", err)
	}

	ensureContents(t, zipFilePath, map[string][]byte{
		"test-dir/test-dir2/file1.txt": []byte("This is file 1"),
		"test-dir/test-file.txt":       []byte("This is test content"),
	})
}

func Keys(m map[string][]byte) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func assertStrings(t *testing.T, message string, want []string, have []string) {
	t.Helper()
	want_i := 0
	have_i := 0
	want_prefix := "  - "
	have_prefix := "  + "
	diff := ""
	for {
		if want_i >= len(want) {
			for ; have_i < len(have); have_i += 1 {
				diff += have_prefix + have[have_i] + "\n"
			}
			break
		}
		if have_i >= len(have) {
			for ; want_i < len(want); want_i += 1 {
				diff += want_prefix + want[want_i] + "\n"
			}
			break
		}
		have_c := have[have_i]
		want_c := want[want_i]
		if want_c < have_c {
			diff += want_prefix + want_c + "\n"
			want_i += 1
			continue
		}
		if have_c < want_c {
			diff += have_prefix + have_c + "\n"
			have_i += 1
			continue
		}
		want_i += 1
		have_i += 1
	}
	if diff != "" {
		t.Fatalf("%s\n%s", message, diff)
	}
}

func assertSameContents(t *testing.T, have map[string][]byte, wants map[string][]byte) {
	assertStrings(t, "file name mismatch:", Keys(wants), Keys(have))
	mismatch := ""
	for name, wantsBytes := range wants {
		hasBytes := have[name]
		wantsStr := string(wantsBytes)
		hasStr := string(hasBytes)
		if wantsStr != hasStr {
			mismatch += fmt.Sprintf("%s:\n  [wants]\n  \n%s\n[has]\n%s", name, wantsStr, hasStr)
		}
	}
	if mismatch != "" {
		t.Fatalf("file content mismatch:\n%s", mismatch)
	}
}

func ensureContents(t *testing.T, zipfilepath string, wants map[string][]byte) {
	t.Helper()
	r, err := zip.OpenReader(zipfilepath)
	if err != nil {
		t.Fatalf("could not open zip file: %s", err)
	}
	defer r.Close()

	have := make(map[string][]byte)
	for _, cf := range r.File {
		if err != nil {
			t.Errorf("could not open file: %s", err)
		}

		r, err := cf.Open()
		if err != nil {
			t.Errorf("could not open file: %s", err)
		}
		defer r.Close()
		gotContentBytes, err := io.ReadAll(r)
		if err != nil {
			t.Errorf("could not read file: %s", err)
		}
		have[cf.Name] = gotContentBytes
	}
	assertSameContents(t, have, wants)
}

func ensureFileMode(t *testing.T, zipfilepath string, outputFileMode string) {
	t.Helper()
	r, err := zip.OpenReader(zipfilepath)
	if err != nil {
		t.Fatalf("could not open zip file: %s", err)
	}
	defer r.Close()

	filemode, err := strconv.ParseUint(outputFileMode, 0, 32)
	if err != nil {
		t.Fatalf("error parsing outputFileMode value: %s", outputFileMode)
	}
	var osfilemode = os.FileMode(filemode)

	for _, cf := range r.File {
		if cf.FileInfo().IsDir() {
			continue
		}

		if cf.Mode() != osfilemode {
			t.Fatalf("Expected filemode \"%s\" but was \"%s\"", osfilemode, cf.Mode())
		}
	}
}
