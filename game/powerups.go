package game

import (
	"gosnake/internal/util"
	"math/rand"
	"time"
)

type PowerUpType int

const (
	SpeedUp PowerUpType = iota - 4
	SlowDown
	GhostMode // Pass through walls temporarily
	ExtraLength
	DoublePoints
)

type PowerUp struct {
	Type     PowerUpType
	Position util.Position
	Duration time.Duration
	Active   bool
	EndTime  time.Time
}

func (g *Game) spawnPowerUp() {
	if g.config.Mode != util.PowerUps {
		return
	}

	if rand.Float32() < 0.2 {
		x, y := g.getRandomEmptyPosition()
		powerType := PowerUpType(rand.Intn(5) - 4)
		g.board[x][y] = int(powerType)
	}
}

func (g *Game) activatePowerUp(typ PowerUpType) {
	powerUp := &PowerUp{
		Type:    typ,
		Active:  true,
		EndTime: time.Now().Add(10 * time.Second),
	}

	switch typ {
	case SpeedUp:
		g.config.Speed = g.config.Speed / 2
	case SlowDown:
		g.config.Speed = g.config.Speed * 2
	case GhostMode:
		g.ghostMode = true
	case ExtraLength:
		g.snake.Length += 2
	case DoublePoints:
		g.pointMultiplier = 2
	}

	g.activePowerUps = append(g.activePowerUps, powerUp)

	g.board[g.snake.Headx][g.snake.Heady] = 0
}

func (g *Game) updatePowerUps() {
	for i := len(g.activePowerUps) - 1; i >= 0; i-- {
		if time.Now().After(g.activePowerUps[i].EndTime) {
			switch g.activePowerUps[i].Type {
			case SpeedUp:
				g.config.Speed *= 2
			case SlowDown:
				g.config.Speed /= 2
			case GhostMode:
				g.ghostMode = false
			case ExtraLength:
				g.snake.Length -= 2
			case DoublePoints:
				g.pointMultiplier = 1
			}
			g.activePowerUps = append(g.activePowerUps[:i], g.activePowerUps[i+1:]...)
		}
	}
}
