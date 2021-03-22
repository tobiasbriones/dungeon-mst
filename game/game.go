/*
 * Copyright (c) 2021 Tobias Briones. All rights reserved.
 */

package game

import (
	"dungeon-mst/ai"
	"dungeon-mst/model"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math/rand"
	"strconv"
	"time"
)

const (
	screenWidth  = 1280
	screenHeight = 720
)

var (
	bgImage  *ebiten.Image
	dungeons []*model.Dungeon
	paths    []*model.Path
)

type Game struct {
	arena       *Arena
	count       int
	legendImage *ebiten.Image
}

func (g *Game) Update() error {
	g.count++

	g.arena.Update(setCurrentDungeonAndPaths)

	// Generate random dungeons
	if g.count%5 == 0 {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			reset()
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(bgImage, nil)

	for _, dungeon := range dungeons {
		dungeon.DrawBarrier(screen)
	}
	for _, path := range paths {
		path.Draw(screen)
	}
	for _, dungeon := range dungeons {
		dungeon.Draw(screen)
	}

	// Draw legend image
	screen.DrawImage(g.legendImage, nil)

	// Draw remote players
	g.arena.Draw(screen)
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) onCharacterMotion(move int) {
	name := g.arena.GetPlayerName()
	println(name + " " + strconv.Itoa(move))
}

func Run() {
	game := newGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Dungeon MST")

	sendFakeInputs(game.arena)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func newGame() Game {
	arena := NewArena()
	legendImage := loadLegendImage()
	game := Game{
		arena:       &arena,
		legendImage: legendImage,
	}

	game.arena.SetOnCharacterMotion(game.onCharacterMotion)
	return game
}

func getSize() model.Dimension {
	return model.NewDimension(screenWidth, screenHeight)
}

func init() {
	loadBg()
	dungeons = ai.GenerateDungeons(getSize())
	//dungeons = genSomeDungeons()

	//genSomeNeighbors(dungeons)
	paths = ai.GetPaths(dungeons)
}

func loadBg() {
	rand.Seed(time.Now().UnixNano())
	bgNumber := rand.Intn(3) + 1
	bgName := "bg_" + strconv.Itoa(bgNumber) + ".png"
	bgImg, _, err := ebitenutil.NewImageFromFile("./assets/" + bgName)

	if err != nil {
		log.Fatal(err)
	}
	bgImage = bgImg
}

func loadLegendImage() *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile("./assets/keyboard_legend.png")

	if err != nil {
		log.Fatal(err)
	}
	return img
}

func setCurrentDungeonAndPaths(runner *model.Runner) {
	var currentDungeon *model.Dungeon = nil
	var currentPaths []*model.Path

	for _, dungeon := range dungeons {
		if dungeon.InBounds(&runner.Rect) {
			currentDungeon = dungeon
			break
		}
	}
	for _, path := range paths {
		if path.InBounds(&runner.Rect) {
			currentPaths = append(currentPaths, path)
		}
	}
	runner.SetCurrentDungeon(currentDungeon)
	runner.SetCurrentPaths(currentPaths)

	if runner.IsOutSide() {
		runner.SetDungeon(dungeons[0])
	}
}

func reset() {
	dungeons = ai.GenerateDungeons(getSize())
	paths = ai.GetPaths(dungeons)
}

func sendFakeInputs(a *Arena) {
	ticker := time.NewTicker(50 * time.Millisecond)

	go func() {
		for range ticker.C {
			fake := randInput()
			a.PushRemotePlayerInput("remote", fake)
		}
	}()
}

func randInput() int {
	return rand.Intn(4)
}
