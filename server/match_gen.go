// Copyright (c) 2021 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/dungeon-mst

package main

import (
	"dungeon-mst/core/geo"
	"dungeon-mst/dungeon"
	"dungeon-mst/game/asset"
	"dungeon-mst/mst"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

func NewRandomMatch() *dungeon.Match {
	dimension := geo.NewDimension(screenWidth, screenHeight)
	return mst.NewRandomMatch(dimension, asset.DiamondSize())
}
