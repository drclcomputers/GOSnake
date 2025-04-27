// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package game

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gosnake/internal/util"

	"github.com/eiannone/keyboard"
)

type Game struct {
	State     util.GameState
	inputChan chan keyboard.KeyEvent
	PowerMgr  util.GamePowerMgr
	sound     *SoundManager
}

func NewGame(Config *util.GameConfig) *Game {
	return &Game{
		State: util.GameState{
			Config:      Config,
			Snake:       util.NewSnake(),
			Board:       util.InitializeBoard(Config.TermWidth, Config.TermHeight),
			Score:       0,
			ExitGame:    false,
			ExitCode:    0,
			PauseGame:   false,
			RelaxedMode: false,
		},

		inputChan: make(chan keyboard.KeyEvent, 10),

		PowerMgr: util.GamePowerMgr{
			GhostMode:       false,
			PointMultiplier: 1,
			ActivePowerUps:  make([]*util.PowerUp, 0),
		},

		sound: NewSoundManager(false),
	}
}

func killSig() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		util.ShowCursor()
		keyboard.Close()
		os.Exit(0)
	}()
}

func (g *Game) placeFood() {
	g.clearEatenApples()

	x, y := g.getRandomEmptyPosition()
	g.State.Board[x][y] = -1
}

func (g *Game) clearEatenApples() {
	for x := 0; x < g.State.Config.TermHeight; x++ {
		for y := 0; y < g.State.Config.TermWidth; y++ {
			if g.State.Board[x][y] == -1 {
				g.State.Board[x][y] = 0
			}
		}
	}
}

func (g *Game) getRandomEmptyPosition() (int, int) {
	for {
		x := rand.Intn(g.State.Config.TermHeight)
		y := rand.Intn(g.State.Config.TermWidth)
		if g.State.Board[x][y] == 0 {
			return x, y
		}
	}
}

func (g *Game) runGameLoop() {
	ticker := time.NewTicker(g.State.Config.Speed)
	defer ticker.Stop()

	renderer := NewRenderer(g.State.Config)
	renderer.Render(g)

	for !g.State.ExitGame {
		g.detectPause()

		select {
		case event := <-g.inputChan:
			g.handleInput(event)
			renderer.Render(g)
		case <-ticker.C:
			g.update()
			if !g.State.ExitGame {
				renderer.Render(g)
			}
		}
	}

	go g.sound.PlayGameOver()

	fmt.Println("\nPress any key to continue...")

	g.writeHighScores()
}
