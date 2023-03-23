// Copyright (c) 2023 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/dungeon-mst

package dungeon

import (
	"dungeon-mst/dungeon"
	"dungeon-mst/game/graphic"
	graphicdungeon "dungeon-mst/game/graphic/dungeon"
	"github.com/hajimehoshi/ebiten/v2"
)

type Dungeon struct {
	*dungeon.Dungeon
	bg      graphic.Draw
	barrier graphic.Draw
}

// Draw make sure to call DrawBarrier first!, it partially draws the dungeon.
func (d Dungeon) Draw(screen *ebiten.Image) {
	d.bg.Draw(screen)
}

// DrawBarrier draws only the Barrier for this Dungeon for correctness issues
// (path being overwritten by barrier).
// It needs to be called before the general call to Draw.
func (d Dungeon) DrawBarrier(screen *ebiten.Image) {
	d.barrier.Draw(screen)
}

func NewDungeonFrom(
	dungeon *dungeon.Dungeon,
	gs graphicdungeon.Graphics,
) *Dungeon {
	return &Dungeon{
		Dungeon: dungeon,
		bg: graphicdungeon.NewDungeonBackgroundDrawing(
			*gs.DungeonBackgroundGraphics,
			dungeon.Rect(),
		),
		barrier: graphicdungeon.NewBrickDrawing(*gs.BrickGraphics, dungeon.Barrier()),
	}
}