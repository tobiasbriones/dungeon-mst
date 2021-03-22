/*
 * Copyright (c) 2021 Tobias Briones. All rights reserved.
 */

package model

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
)

var (
	mplusNormalFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

type Player struct {
	name           string
	character      *Runner
	motionListener MotionListener
}

func (p *Player) GetName() string {
	return p.name
}

func (p *Player) GetCharacter() *Runner {
	return p.character
}

func (p *Player) PushInput(value int) {
	p.character.PushInput(value)
}

func (p *Player) SetMotionListener(value MotionListener) {
	p.motionListener = value
}

func (p *Player) Update() {
	runner := p.character

	if p.motionListener != nil && len(runner.inputs) > 0 {
		p.motionListener(runner.inputs)
	}
	runner.Update()
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.character.Draw(screen)
	p.drawName(screen)
}

func (p *Player) drawName(screen *ebiten.Image) {
	name := p.name
	character := p.character
	x := character.Rect.Left()
	y := character.Rect.Top()
	text.Draw(screen, name, mplusNormalFont, x, y, color.Black)
}

func NewPlayer(name string) Player {
	character := NewRunner()

	return Player{
		name:           name,
		character:      &character,
		motionListener: nil,
	}
}

type MotionListener func([]int)
