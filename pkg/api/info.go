/*
	Copyright 2020 The pdfcpu Authors.

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

	"github.com/alexzh7/pdfcpu/pkg/pdfcpu"
	"github.com/alexzh7/pdfcpu/pkg/pdfcpu/model"
)

// Info returns information about rs.
func Info(rs io.ReadSeeker, selectedPages []string, conf *model.Configuration) ([]string, error) {
	if conf == nil {
		conf = model.NewDefaultConfiguration()
	} else {
		// Validation loads infodict.
		conf.ValidationMode = model.ValidationRelaxed
	}
	ctx, _, _, err := readAndValidate(rs, conf, time.Now())
	if err != nil {
		return nil, err
	}
	if err := ctx.EnsurePageCount(); err != nil {
		return nil, err
	}
	pages, err := PagesForPageSelection(ctx.PageCount, selectedPages, false)
	if err != nil {
		return nil, err
	}
	if err := pdfcpu.DetectWatermarks(ctx); err != nil {
		return nil, err
	}
	return pdfcpu.InfoDigest(ctx, pages)
}

// InfoFile returns information about inFile.
func InfoFile(inFile string, selectedPages []string, conf *model.Configuration) ([]string, error) {
	f, err := os.Open(inFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Info(f, selectedPages, conf)
}