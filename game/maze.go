// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package game

import "gosnake/internal/util"

func (g *Game) generateMaze() {
	g.State.Config.Obstacles = make([]util.Position, 0)

	numObstacles := (g.State.Config.TermWidth * g.State.Config.TermHeight) / 10
	for i := 0; i < numObstacles; i++ {
		x, y := g.getRandomEmptyPosition()
		if x != g.State.Snake.Headx || y != g.State.Snake.Heady {
			obstacle := util.Position{X: x, Y: y}
			g.State.Config.Obstacles = append(g.State.Config.Obstacles, obstacle)
			g.State.Board[x][y] = 999
		}
	}

	g.ensurePlayableMaze()
}

func (g *Game) ensurePlayableMaze() {
	startX, startY := g.State.Snake.Headx, g.State.Snake.Heady
	foodX, foodY := g.getRandomEmptyPosition()
	g.State.Board[foodX][foodY] = -1

	for x := min(startX, foodX); x <= max(startX, foodX); x++ {
		if g.State.Board[x][startY] == 999 {
			g.State.Board[x][startY] = 0
		}
	}
	for y := min(startY, foodY); y <= max(startY, foodY); y++ {
		if g.State.Board[foodX][y] == 999 {
			g.State.Board[foodX][y] = 0
		}
	}
}
