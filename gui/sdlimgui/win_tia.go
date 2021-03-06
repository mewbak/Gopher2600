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
	"gopher2600/hardware/tia/video"

	"github.com/inkyblackness/imgui-go/v2"
)

const winTIATitle = "TIA"

type winTIA struct {
	windowManagement
	img          *SdlImgui
	popupPalette *popupPalette

	strobe int32

	// widget dimensions
	byteDim                     imgui.Vec2
	hmoveSliderWidth            float32
	ballSizeComboDim            imgui.Vec2
	playerSizeAndCopiesComboDim imgui.Vec2
	missileSizeComboDim         imgui.Vec2
	missileCopiesComboDim       imgui.Vec2

	// idxPointer is used to indicate which playfield/player gfx bit is being
	// displayed
	idxPointer imgui.PackedColor
}

func newWinTIA(img *SdlImgui) (managedWindow, error) {
	win := &winTIA{
		img:          img,
		popupPalette: newPopupPalette(img),
		strobe:       -1,
	}

	return win, nil
}

func (win *winTIA) init() {
	win.byteDim = imguiGetFrameDim("FF")
	win.hmoveSliderWidth = imgui.FontSize() * 16
	win.ballSizeComboDim = imguiGetFrameDim("", video.BallSizes...)
	win.playerSizeAndCopiesComboDim = imguiGetFrameDim("", video.PlayerSizes...)
	win.missileSizeComboDim = imguiGetFrameDim("", video.MissileSizes...)
	win.missileCopiesComboDim = imguiGetFrameDim("", video.MissileCopies...)
	win.idxPointer = imgui.PackedColorFromVec4(win.img.cols.IdxPointer)
}

func (win *winTIA) destroy() {
}

func (win *winTIA) id() string {
	return winTIATitle
}

// draw is called by service loop
func (win *winTIA) draw() {
	if !win.open {
		return
	}

	imgui.SetNextWindowPosV(imgui.Vec2{X: 9, Y: 510}, imgui.ConditionFirstUseEver, imgui.Vec2{X: 0, Y: 0})
	imgui.SetNextWindowSizeV(imgui.Vec2{X: 558, Y: 201}, imgui.ConditionFirstUseEver)
	imgui.BeginV(winTIATitle, &win.open, 0)

	imgui.BeginTabBar("")
	if imgui.BeginTabItem("Playfield") {
		win.drawPlayfield()
		imgui.EndTabItem()
	}
	if imgui.BeginTabItem("Player 0") {
		win.drawPlayer(0)
		imgui.EndTabItem()
	}
	if imgui.BeginTabItem("Player 1") {
		win.drawPlayer(1)
		imgui.EndTabItem()
	}
	if imgui.BeginTabItem("Missile 0") {
		win.drawMissile(0)
		imgui.EndTabItem()
	}
	if imgui.BeginTabItem("Missile 1") {
		win.drawMissile(1)
		imgui.EndTabItem()
	}
	if imgui.BeginTabItem("Ball") {
		win.drawBall()
		imgui.EndTabItem()
	}
	imgui.EndTabBar()

	imgui.End()

	win.popupPalette.draw()
}
