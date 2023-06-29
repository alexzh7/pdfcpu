/*
Copyright 2023 The pdfcpu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"io"
	"os"
	"time"

	"github.com/alexzh7/pdfcpu/pkg/log"
	"github.com/alexzh7/pdfcpu/pkg/pdfcpu"
	"github.com/alexzh7/pdfcpu/pkg/pdfcpu/model"
	"github.com/pkg/errors"
)

// ResizeFile applies resizeConf for selected pages of rs and writes result to w.
func Resize(rs io.ReadSeeker, w io.Writer, selectedPages []string, resize *model.Resize, conf *model.Configuration) error {
	if rs == nil {
		return errors.New("pdfcpu: Resize: missing rs")
	}
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	}
	conf.Cmd = model.RESIZE

	ctx, _, _, _, err := readValidateAndOptimize(rs, conf, time.Now())
	if err != nil {
		return err
	}
	if err := ctx.EnsurePageCount(); err != nil {
		return err
	}
	pages, err := PagesForPageSelection(ctx.PageCount, selectedPages, true)
	if err != nil {
		return err
	}

	if err = pdfcpu.Resize(ctx, pages, resize); err != nil {
		return err
	}

	if conf.ValidationMode != model.ValidationNone {
		if err = ValidateContext(ctx); err != nil {
			return err
		}
	}

	return WriteContext(ctx, w)
}

// ResizeFile applies resizeConf for selected pages of inFile and writes result to outFile.
func ResizeFile(inFile, outFile string, selectedPages []string, resize *model.Resize, conf *model.Configuration) error {
	log.CLI.Printf("resizing %s\n", inFile)

	tmpFile := inFile + ".tmp"
	if outFile != "" && inFile != outFile {
		tmpFile = outFile
		log.CLI.Printf("writing %s...\n", outFile)
	} else {
		log.CLI.Printf("writing %s...\n", inFile)
	}

	var (
		f1, f2 *os.File
		err    error
	)

	if f1, err = os.Open(inFile); err != nil {
		return err
	}

	if f2, err = os.Create(tmpFile); err != nil {
		f1.Close()
		return err
	}

	defer func() {
		if err != nil {
			f2.Close()
			f1.Close()
			if outFile == "" || inFile == outFile {
				os.Remove(tmpFile)
			}
			return
		}
		if err = f2.Close(); err != nil {
			return
		}
		if err = f1.Close(); err != nil {
			return
		}
		if outFile == "" || inFile == outFile {
			err = os.Rename(tmpFile, inFile)
		}
	}()

	return Resize(f1, f2, selectedPages, resize, conf)
}
