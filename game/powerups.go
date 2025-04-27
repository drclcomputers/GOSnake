// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package game

import (
	"gosnake/internal/util"
	"math/rand"
	"time"
)

func (g *Game) spawnPowerUp() {
	if g.State.Config.Mode != util.PowerUps {
		return
	}

	if rand.Float32() < 0.2 {
		x, y := g.getRandomEmptyPosition()
		powerType := util.PowerUpType(rand.Intn(5) - 4)
		g.State.Board[x][y] = int(powerType)
	}
}

func (g *Game) activatePowerUp(typ util.PowerUpType) {
	powerUp := &util.PowerUp{
		Type:    typ,
		Active:  true,
		EndTime: time.Now().Add(10 * time.Second),
	}

	switch typ {
	case util.SpeedUp:
		g.State.Config.Speed = g.State.Config.Speed / 2
	case util.SlowDown:
		g.State.Config.Speed = g.State.Config.Speed * 2
	case util.GhostMode:
		g.PowerMgr.GhostMode = true
	case util.ExtraLength:
		g.State.Snake.Length += 2
	case util.DoublePoints:
		g.PowerMgr.PointMultiplier = 2
	}

	g.PowerMgr.ActivePowerUps = append(g.PowerMgr.ActivePowerUps, powerUp)

	g.State.Board[g.State.Snake.Headx][g.State.Snake.Heady] = 0
}

func (g *Game) updatePowerUps() {
	for i := len(g.PowerMgr.ActivePowerUps) - 1; i >= 0; i-- {
		if time.Now().After(g.PowerMgr.ActivePowerUps[i].EndTime) {
			switch g.PowerMgr.ActivePowerUps[i].Type {
			case util.SpeedUp:
				g.State.Config.Speed *= 2
			case util.SlowDown:
				g.State.Config.Speed /= 2
			case util.GhostMode:
				g.PowerMgr.GhostMode = false
			case util.ExtraLength:
				g.State.Snake.Length -= 2
			case util.DoublePoints:
				g.PowerMgr.PointMultiplier = 1
			}
			g.PowerMgr.ActivePowerUps = append(g.PowerMgr.ActivePowerUps[:i], g.PowerMgr.ActivePowerUps[i+1:]...)
		}
	}
}
