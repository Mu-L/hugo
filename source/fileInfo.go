// Copyright 2021 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package source

import (
	"path/filepath"
	"sync"

	"github.com/bep/gitmap"
	"github.com/gohugoio/hugo/common/hashing"
	"github.com/gohugoio/hugo/common/hugo"
	"github.com/gohugoio/hugo/common/paths"
	"github.com/gohugoio/hugo/hugofs/files"

	"github.com/gohugoio/hugo/common/hugio"

	"github.com/gohugoio/hugo/hugofs"
)

// File describes a source file.
type File struct {
	fim hugofs.FileMetaInfo

	uniqueID string
	lazyInit sync.Once
}

// IsContentAdapter returns whether the file represents a content adapter.
// This means that there may be more than one Page associated with this file.
func (fi *File) IsContentAdapter() bool {
	return fi.fim.Meta().PathInfo.IsContentData()
}

// Filename returns a file's absolute path and filename on disk.
func (fi *File) Filename() string { return fi.fim.Meta().Filename }

// Path gets the relative path including file name and extension.  The directory
// is relative to the content root.
func (fi *File) Path() string { return filepath.Join(fi.p().Dir()[1:], fi.p().Name()) }

// Dir gets the name of the directory that contains this file.  The directory is
// relative to the content root.
func (fi *File) Dir() string {
	return fi.pathToDir(fi.p().Dir())
}

// Ext returns a file's extension without the leading period (e.g. "md").
func (fi *File) Ext() string { return fi.p().Ext() }

// Lang returns a file's language (e.g. "sv").
// Deprecated: Use .Page.Language.Lang instead.
func (fi *File) Lang() string {
	hugo.Deprecate(".Page.File.Lang", "Use .Page.Language.Lang instead.", "v0.123.0")
	return fi.fim.Meta().Lang
}

// LogicalName returns a file's name and extension (e.g. "page.sv.md").
func (fi *File) LogicalName() string {
	return fi.p().Name()
}

// BaseFileName returns a file's name without extension (e.g. "page.sv").
func (fi *File) BaseFileName() string {
	return fi.p().NameNoExt()
}

// TranslationBaseName returns a file's translation base name without the
// language segment (e.g. "page").
func (fi *File) TranslationBaseName() string { return fi.p().NameNoIdentifier() }

// ContentBaseName is a either TranslationBaseName or name of containing folder
// if file is a bundle.
func (fi *File) ContentBaseName() string {
	return fi.p().BaseNameNoIdentifier()
}

// Section returns a file's section.
func (fi *File) Section() string {
	return fi.p().Section()
}

// UniqueID returns a file's unique, MD5 hash identifier.
func (fi *File) UniqueID() string {
	fi.init()
	return fi.uniqueID
}

// FileInfo returns a file's underlying os.FileInfo.
func (fi *File) FileInfo() hugofs.FileMetaInfo { return fi.fim }

func (fi *File) String() string { return fi.BaseFileName() }

// Open implements ReadableFile.
func (fi *File) Open() (hugio.ReadSeekCloser, error) {
	f, err := fi.fim.Meta().Open()

	return f, err
}

func (fi *File) IsZero() bool {
	return fi == nil
}

// We create a lot of these FileInfo objects, but there are parts of it used only
// in some cases that is slightly expensive to construct.
func (fi *File) init() {
	fi.lazyInit.Do(func() {
		fi.uniqueID = hashing.MD5FromStringHexEncoded(filepath.ToSlash(fi.Path()))
	})
}

func (fi *File) pathToDir(s string) string {
	if s == "" {
		return s
	}
	return filepath.FromSlash(s[1:] + "/")
}

func (fi *File) p() *paths.Path {
	return fi.fim.Meta().PathInfo.Unnormalized()
}

var contentPathParser = &paths.PathParser{
	IsContentExt: func(ext string) bool {
		return true
	},
}

// Used in tests.
func NewContentFileInfoFrom(path, filename string) *File {
	meta := &hugofs.FileMeta{
		Filename: filename,
		PathInfo: contentPathParser.Parse(files.ComponentFolderContent, filepath.ToSlash(path)),
	}

	return NewFileInfo(hugofs.NewFileMetaInfo(nil, meta))
}

func NewFileInfo(fi hugofs.FileMetaInfo) *File {
	return &File{
		fim: fi,
	}
}

// GitInfo provides information about a version controlled source file.
type GitInfo = gitmap.GitInfo
