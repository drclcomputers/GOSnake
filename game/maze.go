package game

import "gosnake/internal/util"

func (g *Game) generateMaze() {
	g.config.Obstacles = make([]util.Position, 0)

	numObstacles := (g.config.TermWidth * g.config.TermHeight) / 10
	for i := 0; i < numObstacles; i++ {
		x, y := g.getRandomEmptyPosition()
		if x != g.snake.Headx || y != g.snake.Heady {
			obstacle := util.Position{X: x, Y: y}
			g.config.Obstacles = append(g.config.Obstacles, obstacle)
			g.board[x][y] = 999
		}
	}

	g.ensurePlayableMaze()
}

func (g *Game) ensurePlayableMaze() {
	startX, startY := g.snake.Headx, g.snake.Heady
	foodX, foodY := g.getRandomEmptyPosition()
	g.board[foodX][foodY] = -1

	for x := min(startX, foodX); x <= max(startX, foodX); x++ {
		if g.board[x][startY] == 999 {
			g.board[x][startY] = 0
		}
	}
	for y := min(startY, foodY); y <= max(startY, foodY); y++ {
		if g.board[foodX][y] == 999 {
			g.board[foodX][y] = 0
		}
	}
}
