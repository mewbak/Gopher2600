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

package sdlimgui

import (
	"fmt"
	"gopher2600/hardware/memory/memorymap"
	"strconv"

	"github.com/inkyblackness/imgui-go/v2"
)

const winRAMTitle = "RAM"

type winRAM struct {
	windowManagement
	img *SdlImgui

	// widget dimensions
	editWidth float32
	spacing   imgui.Vec2
}

func newWinRAM(img *SdlImgui) (managedWindow, error) {
	win := &winRAM{
		img: img,
	}

	return win, nil
}

func (win *winRAM) init() {
	win.editWidth = minFrameDimension("FF").X
	win.spacing = imgui.Vec2{}
}

func (win *winRAM) destroy() {
}

func (win *winRAM) id() string {
	return winRAMTitle
}

// draw is called by service loop
func (win *winRAM) draw() {
	if !win.open {
		return
	}

	imgui.SetNextWindowPosV(imgui.Vec2{883, 35}, imgui.ConditionFirstUseEver, imgui.Vec2{0, 0})
	imgui.BeginV(winRAMTitle, &win.open, imgui.WindowFlagsAlwaysAutoResize)

	n := 0
	imgui.PushItemWidth(win.editWidth)
	imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, win.spacing)
	for i := memorymap.OriginRAM; i <= memorymap.MemtopRAM; i++ {
		win.drawEditByte(i)
		n++
		if n%16 != 0 {
			imgui.SameLine()
		}
	}
	imgui.PopStyleVar()
	imgui.PopItemWidth()

	imgui.End()
}

func (win *winRAM) drawEditByte(addr uint16) {
	d, _ := win.img.vcs.Mem.RAM.Peek(addr)
	s := fmt.Sprintf("%02x", d)

	cb := func(d imgui.InputTextCallbackData) int32 {
		b := string(d.Buffer())
		// restrict length of input to two characters. note that restriction to
		// hexadecimal characters is handled by imgui's CharsHexadecimal flag
		// given to InputTextV()
		if len(b) > 2 {
			d.DeleteBytes(0, len(b))
			b = b[:2]
			d.InsertBytes(0, []byte(b))
			d.MarkBufferModified()
		}

		return 0
	}

	// flags used with InputTextV()
	flags := imgui.InputTextFlagsCharsHexadecimal |
		imgui.InputTextFlagsCallbackAlways |
		imgui.InputTextFlagsAutoSelectAll

	// if emulator is not paused, the values entered in the TextInput box will
	// be loaded into the register immediately and not just when the enter
	// key is pressed.
	if !win.img.paused {
		flags |= imgui.InputTextFlagsEnterReturnsTrue
	}

	if imgui.InputTextV(fmt.Sprintf("##%d", addr), &s, flags, cb) {
		if v, err := strconv.ParseUint(s, 8, 8); err == nil {
			win.img.vcs.Mem.RAM.Poke(addr, uint8(v))
		}
	}
}
