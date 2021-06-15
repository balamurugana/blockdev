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
	"github.com/balamurugana/blockdev/gpt"
	"github.com/balamurugana/blockdev/mbr"
	"github.com/balamurugana/blockdev/parttable"
	"errors"
	"io"
	"os"
	"reflect"
	"testing"
)

type testPartTable struct {
	uuid       string
	partType   string
	partitions map[int]*parttable.Partition
}

func (tpt *testPartTable) equal(pt parttable.PartTable) bool {
	if pt == nil {
		return false
	}

	if tpt.uuid != pt.UUID() {
		return false
	}

	if tpt.partType != pt.Type() {
		return false
	}

	return reflect.DeepEqual(tpt.partitions, pt.Partitions())
}

func TestMBRProbe(t *testing.T) {
	testCase1Result := &testPartTable{"", "msdos", map[int]*parttable.Partition{}}
	testCase2Result := &testPartTable{
		"",
		"msdos",
		map[int]*parttable.Partition{
			1: &parttable.Partition{Number: 1, Type: parttable.Primary},
			2: &parttable.Partition{Number: 2, Type: parttable.Extended},
			5: &parttable.Partition{Number: 5, Type: parttable.Logical},
			6: &parttable.Partition{Number: 6, Type: parttable.Logical},
		},
	}
	testCase3Result := &testPartTable{
		"",
		"msdos",
		map[int]*parttable.Partition{
			1: &parttable.Partition{Number: 1, Type: parttable.Primary},
			2: &parttable.Partition{Number: 2, Type: parttable.Primary},
			3: &parttable.Partition{Number: 3, Type: parttable.Primary},
			4: &parttable.Partition{Number: 4, Type: parttable.Primary},
		},
	}

	testCases := []struct {
		testDataFile string
		result       *testPartTable
		err          error
	}{
		{"msdos.empty-parts.testdata", testCase1Result, nil},
		{"msdos.logical-partitions.testdata", testCase2Result, nil},
		{"msdos.only-primary-partitions.testdata", testCase3Result, nil},
		{"gpt.testdata", nil, mbr.ErrGPTProtectiveMBR},
		{"zero.testdata", nil, parttable.ErrPartTableNotFound},
	}

	for i, testCase := range testCases {
		devFile, err := os.Open(testCase.testDataFile)
		if err != nil {
			t.Fatalf("case %v: %v: %v", i+1, testCase.testDataFile, err)
		}
		defer devFile.Close()

		result, err := mbr.Probe(devFile)
		if !errors.Is(err, testCase.err) {
			t.Fatalf("case %v: err: expected: %v, got: %v", i+1, testCase.err, err)
		}
		if testCase.result != nil {
			if !testCase.result.equal(result) {
				t.Fatalf("case %v: result: expected: %v, got: %v", i+1, testCase.result, result)
			}
		} else if result != nil {
			t.Fatalf("case %v: result: expected: <nil>, got: %v", i+1, result)
		}
	}
}

func TestGPTProbe(t *testing.T) {
	testCase1Result := &testPartTable{
		"6ce102c7-cfc2-4b1c-b658-02ba8cd9f58f",
		"gpt",
		map[int]*parttable.Partition{
			4: &parttable.Partition{Number: 4, UUID: "8a7d885f-88ba-4734-bbc7-90881480a5a6", Type: parttable.Primary},
			1: &parttable.Partition{Number: 1, UUID: "0d167e49-2c8d-4c6c-ad82-b5e66b6a9eda", Type: parttable.Primary},
			2: &parttable.Partition{Number: 2, UUID: "a183b96b-072c-4236-ae9a-d8adce39859d", Type: parttable.Primary},
			3: &parttable.Partition{Number: 3, UUID: "89fc4f86-1519-47c8-a9f1-11ed504c8f18", Type: parttable.Primary},
		},
	}

	testCases := []struct {
		testDataFile string
		result       *testPartTable
		err          error
	}{
		{"gpt.testdata", testCase1Result, nil},
		{"msdos.empty-parts.testdata", nil, io.EOF},
		{"msdos.logical-partitions.testdata", nil, parttable.ErrPartTableNotFound},
		{"msdos.only-primary-partitions.testdata", nil, io.EOF},
		{"zero.testdata", nil, parttable.ErrPartTableNotFound},
	}

	for i, testCase := range testCases {
		devFile, err := os.Open(testCase.testDataFile)
		if err != nil {
			t.Fatalf("case %v: %v: %v", i+1, testCase.testDataFile, err)
		}
		defer devFile.Close()

		if _, err = devFile.Seek(512, os.SEEK_SET); err != nil {
			t.Fatalf("case %v: %v: %v", i+1, testCase.testDataFile, err)
		}

		result, err := gpt.Probe(devFile)

		if !errors.Is(err, testCase.err) {
			t.Fatalf("case %v: err: expected: %v, got: %v", i+1, testCase.err, err)
		}
		if testCase.result != nil {
			if !testCase.result.equal(result) {
				t.Fatalf("case %v: result: expected: %v, got: %v", i+1, testCase.result, result)
			}
		} else if result != nil {
			t.Fatalf("case %v: result: expected: <nil>, got: %v", i+1, result)
		}
	}
}
