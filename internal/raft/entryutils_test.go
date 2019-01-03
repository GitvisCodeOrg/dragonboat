// Copyright 2017-2019 Lei Ni (nilei81@gmail.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package raft

import (
	"testing"

	pb "github.com/lni/dragonboat/raftpb"
)

func TestNewEntrySlice(t *testing.T) {
	e := make([]pb.Entry, 1, 128)
	n := newEntrySlice(e)
	// TODO
	// should actually check the underly pointer and length to make sure
	// they don't overlap
	if cap(n) == cap(e) {
		t.Errorf("cap is the same")
	}
}

func TestLimitSizeOnEmptyEntryList(t *testing.T) {
	ents := []pb.Entry{}
	limitSize(ents, 0)
}

func TestCheckEntriesToAppend(t *testing.T) {
	checkEntriesToAppend(nil, nil)
	checkEntriesToAppend([]pb.Entry{}, []pb.Entry{})
	checkEntriesToAppend(nil, []pb.Entry{pb.Entry{Index: 101}})
	checkEntriesToAppend([]pb.Entry{pb.Entry{Index: 100}},
		[]pb.Entry{pb.Entry{Index: 101}})
	checkEntriesToAppend([]pb.Entry{pb.Entry{Index: 100, Term: 90}},
		[]pb.Entry{pb.Entry{Index: 101, Term: 100}})
}

func TestCheckEntriesToAppendWillPanicWhenIndexHasHole(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
		t.Errorf("not panic")
	}()
	checkEntriesToAppend([]pb.Entry{pb.Entry{Index: 100}}, []pb.Entry{pb.Entry{Index: 102}})
}

func TestCheckEntriesToAppendWillPanicWhenTermMovesBack(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
		t.Errorf("not panic")
	}()
	checkEntriesToAppend([]pb.Entry{pb.Entry{Index: 100, Term: 100}},
		[]pb.Entry{pb.Entry{Index: 101, Term: 99}})
}