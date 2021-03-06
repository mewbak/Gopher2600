// This file is part of Gopher2600.
//
// Gopher2600 is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Gopher2600 is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Gopher2600.  If not, see <https://www.gnu.org/licenses/>.
//
// *** NOTE: all historical versions of this file, as found in any
// git repository, are also covered by the licence, even when this
// notice is not present ***

package television_test

import (
	"gopher2600/television"
	"testing"
)

func TestNewTelevision(t *testing.T) {
	var tv television.Television
	var err error

	tv, err = television.NewTelevision("PAL")
	if tv == nil || err != nil {
		t.Errorf("PAL spec creation failed")
	}

	tv, err = television.NewTelevision("NTSC")
	if tv == nil || err != nil {
		t.Errorf("NTSC spec creation failed")
	}

	tv, err = television.NewTelevision("AUTO")
	if tv == nil || err != nil {
		t.Errorf("AUTO spec creation failed")
	}

	tv, err = television.NewTelevision("FOO")
	if tv != nil || err == nil {
		t.Errorf("'FOO' spec creation unexpectedly succeeded")
	}
}
