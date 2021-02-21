package model

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

const (
	horizontalUnitWidthPx  = 64
	horizontalUnitHeightPx = 12
	verticalUnitWidthPx    = horizontalUnitHeightPx
	verticalUnitHeightPx   = horizontalUnitWidthPx
)

type Dungeon struct {
	rect        Rect
	brickImage  *ebiten.Image
	brickYImage *ebiten.Image
}

func (d *Dungeon) Width() int {
	return d.rect.Right - d.rect.Left
}

func (d *Dungeon) Height() int {
	return d.rect.Bottom - d.rect.Top
}

func (d *Dungeon) Cx() int {
	return d.rect.Left + d.Width()/2
}

func (d *Dungeon) Cy() int {
	return d.rect.Top + d.Height()/2
}

func (d *Dungeon) Overlaps(other *Dungeon, margin int) bool {
	xo := (d.rect.Left-margin) <= (other.rect.Right+margin) &&
		(d.rect.Right+margin) >= (other.rect.Left-margin)

	yo := (d.rect.Top-margin) <= (other.rect.Bottom+margin) &&
		(d.rect.Bottom+margin) >= (other.rect.Top-margin)

	return xo && yo
}

func (d *Dungeon) Intersects(rect *Rect) bool {
	return d.rect.Intersects(rect)
}

func (d *Dungeon) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	sizeX := d.Width() / horizontalUnitWidthPx
	sizeY := d.Height() / verticalUnitHeightPx
	blockWidth := horizontalUnitHeightPx
	op.GeoM.Translate(float64(d.rect.Left), float64(d.rect.Top))

	for i := 0; i < sizeX; i++ {
		screen.DrawImage(d.brickImage, op)
		op.GeoM.Translate(horizontalUnitWidthPx, 0)
	}

	op.GeoM.Reset()
	op.GeoM.Translate(float64(d.rect.Left), float64(d.rect.Bottom-blockWidth))
	for i := 0; i < sizeX; i++ {
		screen.DrawImage(d.brickImage, op)
		op.GeoM.Translate(horizontalUnitWidthPx, 0)
	}

	op.GeoM.Reset()
	op.GeoM.Translate(float64(d.rect.Left), float64(d.rect.Top))
	for i := 0; i < sizeY; i++ {
		screen.DrawImage(d.brickYImage, op)
		op.GeoM.Translate(0, verticalUnitHeightPx)
	}

	op.GeoM.Reset()
	op.GeoM.Translate(float64(d.rect.Right-blockWidth), float64(d.rect.Top))
	for i := 0; i < sizeY; i++ {
		screen.DrawImage(d.brickYImage, op)
		op.GeoM.Translate(0, verticalUnitHeightPx)
	}
}

func NewDungeon(p0 Point, sizeX int, sizeY int) Dungeon {
	brickImg, _, err := ebitenutil.NewImageFromFile("./assets/brick.png")

	if err != nil {
		log.Fatal(err)
	}
	brickYImg, _, err := ebitenutil.NewImageFromFile("./assets/brick_y.png")

	if err != nil {
		log.Fatal(err)
	}
	rect := Rect{p0.X, p0.Y, p0.X + sizeX*horizontalUnitWidthPx, p0.Y + sizeY*verticalUnitHeightPx}
	return Dungeon{rect, brickImg, brickYImg}
}

func GetDungeonHorizontalUnitSize() Dimension {
	return Dimension{
		Width:  horizontalUnitWidthPx,
		Height: horizontalUnitHeightPx,
	}
}

func GetDungeonVerticalUnitSize() Dimension {
	return Dimension{
		Width:  verticalUnitWidthPx,
		Height: verticalUnitHeightPx,
	}
}
