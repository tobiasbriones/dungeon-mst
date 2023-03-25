// Copyright (c) 2023 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/dungeon-mst

package graphic

import (
	"dungeon-mst/core/geo"
	"github.com/hajimehoshi/ebiten/v2"
)

// Draw defines an object that draws itself on a given canvas like the game
// screen.
type Draw interface {
	Draw(screen *ebiten.Image)
}

// BasicDrawing Defines a simple drawable object that has one Graphic and a
// geo.Rect as position model.
//
// See Draw.
type BasicDrawing struct {
	*Graphic
	*geo.Rect
}

func NewBasicDrawing(graphic *Graphic, rect *geo.Rect) BasicDrawing {
	return BasicDrawing{
		Graphic: graphic,
		Rect:    rect,
	}
}

func (d BasicDrawing) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(d.Left()), float64(d.Top()))
	screen.DrawImage(d.Image, op)
}