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

// Package regression facilitates the regression testing of emulation code. By
// adding test results to a database, the tests can be rerun automatically and
// checked for consistancy.
//
// Currently, two main types of test are supported. First the digest test. This
// test runs a ROM for a set number of frames, saving the video or audio hash
// to the test database.
//
// The second test is the Playback test. This is a slightly more complex test
// that replays user input from a previously recorded session. Recorded
// sessions take video hashes on every input trigger and so will succeed or
// fail if something has changed. The regression test automates the process.
//
// The two tests are useful for different ROMs. The digest type is useful if
// the ROM does something immediately, say an image that is stressful on the
// TIA. The playback type is more useful for real world ROMs (ie. games).
package regression
