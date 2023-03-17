// Copyright (c) 2021 Tobias Briones. All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// This file is part of https://github.com/tobiasbriones/dungeon-mst

package mst

import (
	"dungeon-mst/dungeon"
	"dungeon-mst/geo"
	"math"
	"math/rand"
	"time"
)

const n = 100000
const maxWidthFactor = 8
const maxHeightFactor = 5

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateDungeons(dimension geo.Dimension) []*dungeon.Dungeon {
	var dungeons []*dungeon.Dungeon
	minDim := getMinSize()
	maxDim := getMaxSize()
	xMap := map[int]bool{}
	yMap := map[int]bool{}

	for i := 0; i < n; i++ {
		p := getRandomPoint(dimension, maxDim)
		factor := getRandomFactor()
		w := factor.Width * minDim.Width()
		h := factor.Height * minDim.Width()
		l := p.X() - w/2
		t := p.Y() - h/2
		p0 := geo.NewPoint(l, t)
		rect := geo.NewRect(l, t, l+w, t+h)
		shouldContinue := false

		for _, dungeon := range dungeons {
			if dungeon.Intersects(&rect) {
				shouldContinue = true
				break
			}
		}
		if shouldContinue {
			continue
		}

		// Check if there's a dungeon aligned to this one already
		for i := 0; i <= dungeon.PathWidthPx; i++ {
			if xMap[rect.Left()+i] ||
				xMap[rect.Cx()-dungeon.PathWidthPx/2+i] ||
				xMap[rect.Right()-i] ||
				yMap[rect.Top()+i] ||
				yMap[rect.Cy()-dungeon.PathWidthPx/2+i] ||
				yMap[rect.Bottom()-i] {
				shouldContinue = true
				break
			}
		}
		if shouldContinue {
			continue
		}

		// Update corners
		xMap[rect.Left()] = true
		xMap[rect.Cx()] = true
		xMap[rect.Right()] = true
		yMap[rect.Top()] = true
		yMap[rect.Cy()] = true
		yMap[rect.Bottom()] = true

		// Fill wall widths to avoid paths colliding with walls
		for i := 1; i <= dungeon.PathWidthPx; i++ {
			xMap[rect.Left()+i] = true
			xMap[rect.Right()-i] = true
			xMap[rect.Cx()-dungeon.PathWidthPx/2+i] = true
			yMap[rect.Top()+i] = true
			yMap[rect.Bottom()-i] = true
			yMap[rect.Cy()-dungeon.PathWidthPx/2+i] = true
		}

		// Add the dungeon
		dungeon := dungeon.NewDungeon(p0, factor)
		dungeons = append(dungeons, &dungeon)
	}
	return dungeons
}

func GetPaths(dungeons []*dungeon.Dungeon) []*dungeon.Path {
	var paths []*dungeon.Path
	var tree []*dungeon.Dungeon
	done := map[*dungeon.Dungeon]bool{}

	tree = append(tree, dungeons[0])
	done[dungeons[0]] = true

	for true {
		if len(tree) == len(dungeons) {
			break
		}
		var a *dungeon.Dungeon
		var b *dungeon.Dungeon
		minDistance := 100000

		for _, d1 := range tree {
			p1 := d1.Center()

			for _, d2 := range dungeons {
				if done[d2] {
					continue
				}

				p2 := d2.Center()
				distance := geo.Distance(p1, p2)

				if distance < minDistance {
					minDistance = distance
					a = d1
					b = d2
				}
			}
		}

		if a != nil && b != nil {
			tree = append(tree, b)
			done[b] = true

			path := a.GetPathFor(b)
			//a.AddDoor(path) coming next
			paths = append(paths, path)
		}
	}
	return paths
}

func getMinSize() geo.Dimension {
	size := dungeon.GetDungeonHorizontalUnitSize()
	baseSize := size.Width()
	return geo.NewDimension(baseSize, baseSize)
}

func getMaxSize() geo.Dimension {
	size := dungeon.GetDungeonHorizontalUnitSize()
	baseSize := size.Width()
	return geo.NewDimension(maxWidthFactor*baseSize, maxHeightFactor*baseSize)
}

func getRandomPoint(dimension geo.Dimension, maxDim geo.Dimension) geo.Point {
	cx := maxDim.SemiWidth() + int(float64(dimension.Width()-maxDim.Width())*rand.Float64())
	cy := maxDim.SemiHeight() + int(float64(dimension.Height()-maxDim.Height())*rand.Float64())
	return geo.NewPoint(cx, cy)
}

func getRandomFactor() dungeon.DimensionFactor {
	wFactor := 1 + int(math.Floor(float64(maxWidthFactor)*rand.Float64()))
	hFactor := 1 + int(math.Floor(float64(maxHeightFactor)*rand.Float64()))
	return dungeon.DimensionFactor{Width: wFactor, Height: hFactor}
}