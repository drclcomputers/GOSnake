// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package game

import (
	"fmt"
	"gosnake/internal/util"
	"strings"
)

func (g *Game) moveSnake() {
	newX, newY := g.State.Snake.Headx, g.State.Snake.Heady

	switch g.State.Snake.Direction {
	case 1:
		newX--
	case 2:
		newY++
	case 3:
		newX++
	case 4:
		newY--
	}

	if g.State.Config.Mode == util.NoWalls || g.PowerMgr.GhostMode {
		if newX < 0 {
			newX = g.State.Config.TermHeight - 1
		} else if newX >= g.State.Config.TermHeight {
			newX = 0
		}
		if newY < 0 {
			newY = g.State.Config.TermWidth - 1
		} else if newY >= g.State.Config.TermWidth {
			newY = 0
		}
	}

	g.State.Snake.Headx, g.State.Snake.Heady = newX, newY
}

func (g *Game) checkCollision() int {
	if g.State.Config.Mode != util.NoWalls && !g.PowerMgr.GhostMode {
		if g.State.Snake.Headx < 0 || g.State.Snake.Headx >= g.State.Config.TermHeight ||
			g.State.Snake.Heady < 0 || g.State.Snake.Heady >= g.State.Config.TermWidth {
			return util.CollisionWall
		}
	}

	if g.State.Config.Mode == util.Maze && !g.PowerMgr.GhostMode {
		if g.State.Board[g.State.Snake.Headx][g.State.Snake.Heady] == 999 {
			return util.CollisionWall
		}
	}

	if g.State.Board[g.State.Snake.Headx][g.State.Snake.Heady] > 0 {
		return util.CollisionSelf
	}

	return 0
}

func (g *Game) updateBoard() {
	obstacles := make(map[util.Position]bool)
	if g.State.Config.Mode == util.Maze {
		for _, obs := range g.State.Config.Obstacles {
			obstacles[obs] = true
		}
	}

	switch g.State.Board[g.State.Snake.Headx][g.State.Snake.Heady] {
	case -1: // Food
		g.State.Score += (1 * g.PowerMgr.PointMultiplier)
		g.State.Snake.Length++
		if g.sound != nil {
			go g.sound.PlayFoodEaten()
		}
		g.placeFood()

		if g.State.Config.Mode == util.PowerUps {
			g.spawnPowerUp()
		}
	case -2, -3, -4, -5, -6: // Powerups
		if g.sound != nil {
			go g.sound.PlayPowerUpCollected()
		}
		g.activatePowerUp(util.PowerUpType(g.State.Board[g.State.Snake.Headx][g.State.Snake.Heady]))
	}

	g.updatePowerUps()

	for x := 0; x < g.State.Config.TermHeight; x++ {
		for y := 0; y < g.State.Config.TermWidth; y++ {
			if g.State.Board[x][y] > 0 {
				g.State.Board[x][y]++
			}

			if g.State.Board[x][y] > g.State.Snake.Length {
				g.State.Board[x][y] = 0
			}
		}
	}

	if g.State.Config.Mode == util.Maze {
		for pos := range obstacles {
			g.State.Board[pos.X][pos.Y] = 999
		}
	}

	g.State.Board[g.State.Snake.Headx][g.State.Snake.Heady] = 1
}

func (g *Game) update() {
	g.moveSnake()
	col := g.checkCollision()
	if col != 0 {
		g.State.ExitCode = col
		g.State.ExitGame = true
		StopMusic()
		if col == 1 {
			fmt.Println("\n" + strings.Repeat(" ", g.State.Config.OffsetX-1) + "Game Over! Snake hit a wall!")
		} else if col == 2 {
			fmt.Println("\n" + strings.Repeat(" ", g.State.Config.OffsetX-1) + "Game Over! Snake collided with itself!")
		}
		return
	}
	g.updateBoard()
	if g.State.Score%util.MODSPEED == 0 && g.State.Score > 0 && !g.State.RelaxedMode {
		g.State.Config.Speed -= util.MODSPEED
	}
}

func (g *Game) detectPause() {
	for g.State.PauseGame {
		select {
		case event := <-g.inputChan:
			g.handleInput(event)
		}
	}
}
