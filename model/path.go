/*
 * Copyright (c) 2021 Tobias Briones. All rights reserved.
 */

package model

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

const (
	PathWidthPx = 32
)

type Path struct {
	hLine      Line
	hRect      Rect
	hRectImage *ebiten.Image
	vLine      Line
	vRect      Rect
	vRectImage *ebiten.Image
}

func (p *Path) InBounds(rect *Rect) bool {
	hRect := p.hRect
	vRect := p.vRect
	return hRect.InBounds(rect) || vRect.InBounds(rect)
}

func (p *Path) CanMoveTowards(movement Movement, rect *Rect) bool {
	if !p.InBounds(rect) {
		return true
	}
	if p.hRect.InBounds(rect) {
		return CheckMovement(movement, rect, p.hRect)
	}
	return CheckMovement(movement, rect, p.vRect)
}

func CheckMovement(movement Movement, rect *Rect, host Rect) bool {
	if movement.direction == MoveDirLeft {
		return rect.Left-movement.length > host.Left
	} else if movement.direction == MoveDirTop {
		return rect.Top-movement.length > host.Top
	} else if movement.direction == MoveDirRight {
		return rect.Right+movement.length < host.Right
	} else if movement.direction == MoveDirBottom {
		return rect.Bottom+movement.length < host.Bottom
	}
	return false
}

func (p *Path) Draw(screen *ebiten.Image) {
	p.drawHorizontalLine(screen)
	p.drawVerticalLine(screen)
}

func (p *Path) drawHorizontalLine(screen *ebiten.Image) {
	rect := p.hRect
	x := rect.Left
	y := rect.Top
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(p.hRectImage, op)
}

func (p *Path) drawVerticalLine(screen *ebiten.Image) {
	rect := p.vRect
	x := rect.Left
	y := rect.Top
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(p.vRectImage, op)
}

func NewPath(hl Line, vl Line) Path {
	if hl.IsDegenerate() || vl.IsDegenerate() {
		panic("Lines cannot be degenerate")
	}
	if !(hl.IsHorizontal() && vl.IsVertical()) {
		panic("The first line must be horizontal and the other vertical")
	}
	if !hl.p1.Equals(&vl.p1) && !hl.p1.Equals(&vl.p2) &&
		!hl.p2.Equals(&vl.p1) && !hl.p2.Equals(&vl.p2) {
		panic("One of the horizontal line start/end point must be the beginning of the vertical line")
	}
	if hl.p1.X > hl.p2.X {
		panic("The point 1 of the horizontal line must be the lowest")
	}
	if vl.p1.Y > vl.p2.Y {
		panic("The point 1 of the vertical line must be the lowest")
	}
	sw := PathWidthPx / 2
	hRect := Rect{
		Left:   hl.p1.X,
		Top:    hl.p1.Y - sw,
		Right:  hl.p2.X,
		Bottom: hl.p1.Y + sw,
	}
	vRect := Rect{
		Left:   vl.p1.X - sw,
		Top:    vl.p1.Y,
		Right:  vl.p1.X + sw,
		Bottom: vl.p2.Y,
	}

	hImg := ebiten.NewImage(hRect.Width(), hRect.Height())
	vImg := ebiten.NewImage(vRect.Width(), vRect.Height())

	hImg.Fill(color.Gray{})
	vImg.Fill(color.Gray{})
	return Path{hl, hRect, hImg, vl, vRect, vImg}
}

type Line struct {
	p1 Point
	p2 Point
}

func (l *Line) IsDegenerate() bool {
	return Distance(l.p1, l.p2) == 0
}

func (l *Line) IsHorizontal() bool {
	return l.p1.Y == l.p2.Y
}

func (l *Line) IsVertical() bool {
	return l.p1.X == l.p2.X
}