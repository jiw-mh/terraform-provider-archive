// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package archive

import (
	"fmt"
	"path/filepath"
	"regexp"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZipArchiveFile_Basic(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string
	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccArchiveFileContentConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "2979b8426e585514344e6d78b17ed5d9",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "6d0527096342115cfb55b18f8420b695d91e2622",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha256", "865db2315db933c59bc9268455ddd02fd73c501d9547434c42295cccc6dda4a4",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "hl2yMV25M8WbySaEVd3QL9c8UB2VR0NMQilczMbdpKQ=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha512", "25d2bfbf3bc4daa363433b0167c3a58c9b97a4fcee8a58111fda3d5529b81a4c3527a8f18eadac592722ac10f09db27dd4a9a80ee2dbde4ad6b9a564613d2ce0",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha512", "JdK/vzvE2qNjQzsBZ8OljJuXpPzuilgRH9o9VSm4Gkw1J6jxjq2sWScirBDwnbJ91KmoDuLb3krWuaVkYT0s4A==",
					),
				),
			},
			{
				Config: testAccArchiveFileFileConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "2985d42029a4249a50ef480dfbee3dac",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "653f0bb73c57782db99e5b8bbdda5bab70250240",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha256", "4478a74a621dc9d702d99adc0494f8b8556365a693b3b128e6b203a47b917717",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "RHinSmIdydcC2ZrcBJT4uFVjZaaTs7Eo5rIDpHuRdxc=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha512", "0b545b786080a05a87690fc1a17329967532970e435d80a34627b79d4e73a1c80ce980610b439cf4eda62ae53e4f560115cbbfccfd2650ffe6712656d8db323b",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha512", "C1RbeGCAoFqHaQ/BoXMplnUylw5DXYCjRie3nU5zocgM6YBhC0Oc9O2mKuU+T1YBFcu/zP0mUP/mcSZW2NsyOw==",
					),
				),
			},
			{
				Config: testAccArchiveFileDirConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_md5", "53e1fed54b3a677717dffc50f50bdcf8",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha", "914fc63a70af426f7911567e1360911f23de8d2d",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha256", "a0f6160168c1a44dc2e4c63b312fe862c37ab8327e3de61cdc97253751a33fce",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha256", "oPYWAWjBpE3C5MY7MS/oYsN6uDJ+PeYc3JclN1GjP84=",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_sha512", "0664c096095836ca09f207c465a80f61a05377209ca4e21a9238c7bec10840e42c4476f9c13ffcd1bd1c14da221e8a512e423393437dded5fb2bba13b7160cd7",
					),
					r.TestCheckResourceAttr(
						"data.archive_file.foo", "output_base64sha512", "BmTAlglYNsoJ8gfEZagPYaBTdyCcpOIakjjHvsEIQOQsRHb5wT/80b0cFNoiHopRLkIzk0N93tX7K7oTtxYM1w==",
					),
				),
			},
			{
				Config: testAccArchiveFileDirExcludesConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				Config: testAccArchiveFileDirExcludesGlobConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
			{
				Config: testAccArchiveFileMultiSourceConfig("zip", f),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
				),
			},
		},
	})
}

func TestAccZipArchiveFile_SourceConfigMissing(t *testing.T) {
	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config:      testAccArchiveSourceConfigMissing("zip"),
				ExpectError: regexp.MustCompile(`.*At least one of these attributes must be configured:\n\[source,source_content_filename,source_file,source_dir]`),
			},
		},
	})
}

func TestAccZipArchiveFile_SourceConfigConflicting(t *testing.T) {
	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config:      testAccArchiveSourceConfigConflicting("zip"),
				ExpectError: regexp.MustCompile(`.*Attribute "source_dir" cannot be specified when "source" is specified`),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkFile_Relative verifies that a symlink to a file using a relative path generates an
// archive which includes the file.
func TestAccZipArchiveFile_SymlinkFile_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkFile_Absolute verifies that a symlink to a file using an absolute path generates an
// archive which includes the file.
func TestAccZipArchiveFile_SymlinkFile_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkFileAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkDirectory_Relative verifies that a symlink to a directory using a relative path
// generates an archive which includes the files in the directory.
func TestAccZipArchiveFile_SymlinkDirectory_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"file1.txt": []byte(`This is file 1`),
							"file2.txt": []byte(`This is file 2`),
							"file3.txt": []byte(`This is file 3`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkDirectory_Absolute verifies that a symlink to a directory using an absolute path
// generates an archive which includes the files in the directory.
func TestAccZipArchiveFile_SymlinkDirectory_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithRegularFilesAbs, err := filepath.Abs("test-fixtures/test-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirWithRegularFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"file1.txt": []byte(`This is file 1`),
							"file2.txt": []byte(`This is file 2`),
							"file3.txt": []byte(`This is file 3`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_DirectoryWithSymlinkFile_Relative verifies that a relative path to a directory containing
// a symlink file generates an archive which includes the files in the directory.
func TestAccZipArchiveFile_DirectoryWithSymlinkFile_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-file.txt":    []byte(`This is test content`),
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_DirectoryWithSymlinkFile_Absolute verifies that an absolute path to a directory containing
// a symlink file generates an archive which includes the files in the directory.
func TestAccZipArchiveFile_DirectoryWithSymlinkFile_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithSymlinkFilesAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirWithSymlinkFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-file.txt":    []byte(`This is test content`),
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative verifies that a relative path to a symlink
// file in a symlink directory generates an archive which includes the files in the directory.
func TestAccZipArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute verifies that an absolute path to a symlink
// file in a symlink directory generates an archive which includes the files in the directory.
func TestAccZipArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileInSymlinkDirAbs, err := filepath.Abs("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_file      = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkFileInSymlinkDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_DirectoryWithSymlinkDirectory_Relative verifies that a relative path to a
// directory containing a symlink to a directory generates an archive which includes the directory.
func TestAccZipArchiveFile_DirectoryWithSymlinkDirectory_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-dir"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink-dir/file1.txt": []byte("This is file 1"),
							"test-symlink-dir/file2.txt": []byte("This is file 2"),
							"test-symlink-dir/file3.txt": []byte("This is file 3"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute verifies that an absolute path to a
// directory containing a symlink to a directory generates an archive which includes the directory.
func TestAccZipArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirInRegularDirAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(symlinkDirInRegularDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink-dir/file1.txt": []byte("This is file 1"),
							"test-symlink-dir/file2.txt": []byte("This is file 2"),
							"test-symlink-dir/file3.txt": []byte("This is file 3"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_Multiple_Relative verifies that a relative path to a directory containing multiple
// directories including symlink directories generates an archive which includes the directories and files.
func TestAccZipArchiveFile_Multiple_Relative(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash("test-fixtures"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-dir/test-dir1/file1.txt":                         []byte("This is file 1"),
							"test-dir/test-dir1/file2.txt":                         []byte("This is file 2"),
							"test-dir/test-dir1/file3.txt":                         []byte("This is file 3"),
							"test-dir/test-dir2/file1.txt":                         []byte("This is file 1"),
							"test-dir/test-dir2/file2.txt":                         []byte("This is file 2"),
							"test-dir/test-dir2/file3.txt":                         []byte("This is file 3"),
							"test-dir/test-file.txt":                               []byte("This is test content"),
							"test-dir-with-symlink-dir/test-symlink-dir/file1.txt": []byte("This is file 1"),
							"test-dir-with-symlink-dir/test-symlink-dir/file2.txt": []byte("This is file 2"),
							"test-dir-with-symlink-dir/test-symlink-dir/file3.txt": []byte("This is file 3"),
							"test-dir-with-symlink-file/test-file.txt":             []byte("This is test content"),
							"test-dir-with-symlink-file/test-symlink.txt":          []byte("This is test content"),
							"test-symlink-dir/file1.txt":                           []byte("This is file 1"),
							"test-symlink-dir/file2.txt":                           []byte("This is file 2"),
							"test-symlink-dir/file3.txt":                           []byte("This is file 3"),
							"test-symlink-dir-with-symlink-file/test-file.txt":     []byte("This is test content"),
							"test-symlink-dir-with-symlink-file/test-symlink.txt":  []byte("This is test content"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_Multiple_Absolute verifies that an absolute path to a directory containing multiple
// directories including symlink directories generates an archive which includes the directories and files.
func TestAccZipArchiveFile_Multiple_Absolute(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	multipleDirsAndFilesAbs, err := filepath.Abs("test-fixtures")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type             = "zip"
			 source_dir       = "%s"
			 output_path      = "%s"
			 output_file_mode = "0666"
			}
			`, filepath.ToSlash(multipleDirsAndFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-dir/test-dir1/file1.txt":                         []byte("This is file 1"),
							"test-dir/test-dir1/file2.txt":                         []byte("This is file 2"),
							"test-dir/test-dir1/file3.txt":                         []byte("This is file 3"),
							"test-dir/test-dir2/file1.txt":                         []byte("This is file 1"),
							"test-dir/test-dir2/file2.txt":                         []byte("This is file 2"),
							"test-dir/test-dir2/file3.txt":                         []byte("This is file 3"),
							"test-dir/test-file.txt":                               []byte("This is test content"),
							"test-dir-with-symlink-dir/test-symlink-dir/file1.txt": []byte("This is file 1"),
							"test-dir-with-symlink-dir/test-symlink-dir/file2.txt": []byte("This is file 2"),
							"test-dir-with-symlink-dir/test-symlink-dir/file3.txt": []byte("This is file 3"),
							"test-dir-with-symlink-file/test-file.txt":             []byte("This is test content"),
							"test-dir-with-symlink-file/test-symlink.txt":          []byte("This is test content"),
							"test-symlink-dir/file1.txt":                           []byte("This is file 1"),
							"test-symlink-dir/file2.txt":                           []byte("This is file 2"),
							"test-symlink-dir/file3.txt":                           []byte("This is file 3"),
							"test-symlink-dir-with-symlink-file/test-file.txt":     []byte("This is test content"),
							"test-symlink-dir-with-symlink-file/test-symlink.txt":  []byte("This is test content"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a symlink to a file using a relative
// path generates an archive which includes the file.
func TestAccZipArchiveFile_SymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that a symlink to a file using an absolute
// path generates an archive which includes the file.
func TestAccZipArchiveFile_SymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkFileAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkDirectory_Relative_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestAccZipArchiveFile_SymlinkDirectory_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir"), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: archive has not been\ncreated as it would be empty`),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkDirectory_Absolute_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestAccZipArchiveFile_SymlinkDirectory_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithRegularFilesAbs, err := filepath.Abs("test-fixtures/test-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkDirWithRegularFilesAbs), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: archive has not been\ncreated as it would be empty`),
			},
		},
	})
}

// TestAccZipArchiveFile_DirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a relative path to a
// directory containing a symlink file generates an archive which includes the files in the directory.
func TestAccZipArchiveFile_DirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-file"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-file.txt":    []byte(`This is test content`),
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_DirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that an absolute path to a
// directory containing a symlink file generates an archive which includes the files in the directory.
func TestAccZipArchiveFile_DirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirWithSymlinkFilesAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-file")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkDirWithSymlinkFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-file.txt":    []byte(`This is test content`),
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories verifies that a relative path
// to a symlink file in a symlink directory generates an archive which includes the files in the directory.
func TestAccZipArchiveFile_SymlinkDirectoryWithSymlinkFile_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories verifies that an absolute path
// to a symlink file in a symlink directory generates an archive which includes the files in the directory.
func TestAccZipArchiveFile_SymlinkDirectoryWithSymlinkFile_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkFileInSymlinkDirAbs, err := filepath.Abs("test-fixtures/test-symlink-dir-with-symlink-file/test-symlink.txt")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_file                 = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkFileInSymlinkDirAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-symlink.txt": []byte(`This is test content`),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_DirectoryWithSymlinkDirectory_Relative_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestAccZipArchiveFile_DirectoryWithSymlinkDirectory_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures/test-dir-with-symlink-dir"), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: archive has not been\ncreated as it would be empty`),
			},
		},
	})
}

// TestAccZipArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute_ExcludeSymlinkDirectories verifies that an empty archive
// is generated when trying to archive a directory which only contains a symlink to a directory.
func TestAccZipArchiveFile_IncludeDirectoryWithSymlinkDirectory_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	symlinkDirInRegularDirAbs, err := filepath.Abs("test-fixtures/test-dir-with-symlink-dir")
	if err != nil {
		t.Fatal(err)
	}

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(symlinkDirInRegularDirAbs), filepath.ToSlash(f)),
				ExpectError: regexp.MustCompile(`.*error creating archive: error archiving directory: archive has not been\ncreated as it would be empty`),
			},
		},
	})
}

// TestAccZipArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories verifies that
// symlinked directories are excluded.
func TestAccZipArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash("test-fixtures"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-dir/test-dir1/file1.txt":                []byte("This is file 1"),
							"test-dir/test-dir1/file2.txt":                []byte("This is file 2"),
							"test-dir/test-dir1/file3.txt":                []byte("This is file 3"),
							"test-dir/test-dir2/file1.txt":                []byte("This is file 1"),
							"test-dir/test-dir2/file2.txt":                []byte("This is file 2"),
							"test-dir/test-dir2/file3.txt":                []byte("This is file 3"),
							"test-dir/test-file.txt":                      []byte("This is test content"),
							"test-dir-with-symlink-file/test-file.txt":    []byte("This is test content"),
							"test-dir-with-symlink-file/test-symlink.txt": []byte("This is test content"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestAccZipArchiveFile_Multiple_Relative_ExcludeSymlinkDirectories verifies that
// symlinked directories are excluded.
func TestAccZipArchiveFile_Multiple_Absolute_ExcludeSymlinkDirectories(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	multipleDirsAndFilesAbs, err := filepath.Abs("test-fixtures")
	if err != nil {
		t.Fatal(err)
	}

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source_dir                  = "%s"
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 exclude_symlink_directories = true
			}
			`, filepath.ToSlash(multipleDirsAndFilesAbs), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test-dir/test-dir1/file1.txt":                []byte("This is file 1"),
							"test-dir/test-dir1/file2.txt":                []byte("This is file 2"),
							"test-dir/test-dir1/file3.txt":                []byte("This is file 3"),
							"test-dir/test-dir2/file1.txt":                []byte("This is file 1"),
							"test-dir/test-dir2/file2.txt":                []byte("This is file 2"),
							"test-dir/test-dir2/file3.txt":                []byte("This is file 3"),
							"test-dir/test-file.txt":                      []byte("This is test content"),
							"test-dir-with-symlink-file/test-file.txt":    []byte("This is test content"),
							"test-dir-with-symlink-file/test-symlink.txt": []byte("This is test content"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestZipArchiveFile_Source verifies that source configuration is included
func TestZipArchiveFile_Source(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source {
			    content = "hello world"
				filename = "test.txt"
			 }
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			}
			`, filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test.txt": []byte("hello world"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestZipArchiveFile_Source verifies that source configuration is included
func TestZipArchiveFile_SourceExcluded(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source {
			    content = "hello world"
				filename = "test.txt"
			 }
		     source {
				content = "foo"
				filename = "foo/bar.txt"
			 }
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			 excludes = ["foo/*"]
			}
			`, filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test.txt": []byte("hello world"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestZipArchiveFile_SourcePath verifies that source content as base64 is included
func TestZipArchiveFile_SourcePath(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source {
			    source_path = "%s"
				filename = "test.txt"
			 }
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			}
			`, filepath.ToSlash("test-fixtures/test-dir/test-file.txt"), filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test.txt": []byte("This is test content"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}

// TestZipArchiveFile_ContentBase64 verifies that source content as base64 is included
func TestZipArchiveFile_ContentBase64(t *testing.T) {
	td := t.TempDir()

	f := filepath.Join(td, "zip_file_acc_test.zip")

	var fileSize string

	r.ParallelTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
			data "archive_file" "foo" {
			 type                        = "zip"
			 source {
			    content_base64 = "aGVsbG8K"
				filename = "test.txt"
			 }
			 output_path                 = "%s"
			 output_file_mode            = "0666"
			}
			`, filepath.ToSlash(f)),
				Check: r.ComposeTestCheckFunc(
					testAccArchiveFileSize(f, &fileSize),
					r.TestCheckResourceAttrPtr("data.archive_file.foo", "output_size", &fileSize),
					r.TestCheckResourceAttrWith("data.archive_file.foo", "output_path", func(value string) error {
						ensureContents(t, value, map[string][]byte{
							"test.txt": []byte("hello\n"),
						})
						ensureFileMode(t, value, "0666")
						return nil
					}),
				),
			},
		},
	})
}
