package game

import (
	"gosnake/internal/util"
	"math/rand"
)

func (g *Game) moveSnake() {
	newX, newY := g.snake.Headx, g.snake.Heady

	switch g.snake.Direction {
	case 1:
		newX--
	case 2:
		newY++
	case 3:
		newX++
	case 4:
		newY--
	}

	if g.config.Mode == util.NoWalls || g.ghostMode {
		if newX < 0 {
			newX = g.config.TermHeight - 1
		} else if newX >= g.config.TermHeight {
			newX = 0
		}
		if newY < 0 {
			newY = g.config.TermWidth - 1
		} else if newY >= g.config.TermWidth {
			newY = 0
		}
	}

	g.snake.Headx, g.snake.Heady = newX, newY
}

func (g *Game) checkCollision() int {
	if g.config.Mode != util.NoWalls && !g.ghostMode {
		if g.snake.Headx < 0 || g.snake.Headx >= g.config.TermHeight ||
			g.snake.Heady < 0 || g.snake.Heady >= g.config.TermWidth {
			return 1
		}
	}

	if g.config.Mode == util.Maze && !g.ghostMode {
		if g.board[g.snake.Headx][g.snake.Heady] == 999 {
			return 1
		}
	}

	if g.board[g.snake.Headx][g.snake.Heady] > 0 {
		return 2
	}

	return 0
}

func (g *Game) updateBoard() {
	obstacles := make(map[util.Position]bool)
	if g.config.Mode == util.Maze {
		for _, obs := range g.config.Obstacles {
			obstacles[obs] = true
		}
	}

	switch g.board[g.snake.Headx][g.snake.Heady] {
	case -1: // Food
		g.score += (1 * g.pointMultiplier)
		g.snake.Length++
		if g.sound != nil {
			go g.sound.PlayFoodEaten()
		}
		g.placeFood()

		if g.config.Mode == util.PowerUps {
			g.spawnPowerUp()
		}
	case -2, -3, -4, -5, -6: // Powerups
		if g.sound != nil {
			go g.sound.PlayPowerUpCollected()
		}
		g.activatePowerUp(PowerUpType(g.board[g.snake.Headx][g.snake.Heady]))
	}

	g.updatePowerUps()

	for x := 0; x < g.config.TermHeight; x++ {
		for y := 0; y < g.config.TermWidth; y++ {
			if g.board[x][y] > 0 {
				g.board[x][y]++
			}

			if g.board[x][y] > g.snake.Length {
				g.board[x][y] = 0
			}
		}
	}

	if g.config.Mode == util.Maze {
		for pos := range obstacles {
			g.board[pos.X][pos.Y] = 999
		}
	}

	g.board[g.snake.Headx][g.snake.Heady] = 1
}

func (g *Game) getRandomEmptyPosition() (int, int) {
	for {
		x := rand.Intn(g.config.TermHeight)
		y := rand.Intn(g.config.TermWidth)
		if g.board[x][y] == 0 {
			return x, y
		}
	}
}
