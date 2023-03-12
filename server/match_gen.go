// Copyright (c) 2021 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/dungeon-mst

package main

import (
	"dungeon-mst/game/model"
	"dungeon-mst/geo"
	"dungeon-mst/mst"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

func NewRandomMatch() *model.Match {
	dimension := geo.NewDimension(screenWidth, screenHeight)
	return mst.NewRandomMatch(dimension)
}
