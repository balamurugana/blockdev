/*
 * Library to probe partition table attributes from block device.
 * Copyright 2021 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package blockdev

import (
	"errors"
	"github.com/balamurugana/blockdev/gpt"
	"github.com/balamurugana/blockdev/mbr"
	"github.com/balamurugana/blockdev/parttable"
	"os"
)

// Probe detects and returns partition table in given device filename.
func Probe(filename string) (parttable.PartTable, error) {
	devFile, err := os.OpenFile(filename, os.O_RDONLY, os.ModeDevice)
	if err != nil {
		return nil, err
	}
	defer devFile.Close()

	partTable, err := mbr.Probe(devFile)
	if err == nil {
		return partTable, nil
	}
	if !errors.Is(err, parttable.ErrPartTableNotFound) && !errors.Is(err, mbr.ErrGPTProtectiveMBR) {
		return nil, err
	}

	if _, err = devFile.Seek(512, os.SEEK_SET); err != nil {
		return nil, err
	}
	return gpt.Probe(devFile)
}
