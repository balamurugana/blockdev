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

package parttable

import (
	"errors"
)

var ErrPartTableNotFound = errors.New("partition table not found")

type PartType int

const (
	Primary PartType = iota + 1
	Extended
	Logical
)

func (pt PartType) String() string {
	switch pt {
	case Primary:
		return "primary"
	case Extended:
		return "extended"
	case Logical:
		return "logical"
	default:
		return ""
	}
}

// Partition denotes partition information.
type Partition struct {
	Number int
	UUID   string
	Type   PartType
}

// PartTable denotes partition table.
type PartTable interface {
	UUID() string
	Type() string
	Partitions() map[int]*Partition
}
