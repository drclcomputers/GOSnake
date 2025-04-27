// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package game

import (
	"fmt"
	"gosnake/internal/util"
	"strconv"
	"strings"
	"time"
)

type Renderer struct {
	Config *util.GameConfig
}

func NewRenderer(Config *util.GameConfig) *Renderer {
	return &Renderer{Config: Config}
}

func (r *Renderer) Render(g *Game) {
	var builder strings.Builder

	r.renderTopOffset(&builder)
	r.renderTopAndBottomBorder(&builder, g.PowerMgr.GhostMode)
	r.renderBoard(&builder, g)
	r.renderTopAndBottomBorder(&builder, g.PowerMgr.GhostMode)
	r.renderScore(&builder, g.State.Score)
	if g.State.Config.Mode == util.PowerUps {
		r.renderActiveEffects(&builder, g)
	}

	util.GoAtTopLeft()
	fmt.Print(builder.String())
}

func (r *Renderer) decideColor(builder *strings.Builder, isGhostMode bool) {
	builder.WriteString(util.BLACK)

	if isGhostMode {
		builder.WriteString(util.GREEN)
		return
	}

	switch r.Config.Mode {
	case util.NoWalls:
		builder.WriteString(util.GREEN)
	default:
		builder.WriteString(util.RED)
	}
}

func (r *Renderer) renderBoard(builder *strings.Builder, g *Game) {
	for x := 0; x < g.State.Config.TermHeight; x++ {
		r.decideColor(builder, g.PowerMgr.GhostMode)
		builder.WriteString(strings.Repeat(" ", g.State.Config.OffsetX-1) + g.State.Config.BorderChar)
		builder.WriteString(util.BLACK)

		for y := 0; y < g.State.Config.TermWidth; y++ {
			switch {
			case g.State.Board[x][y] == 0:
				builder.WriteString(g.State.Config.EmptyCell)
			case g.State.Board[x][y] == -1:
				builder.WriteString(g.State.Config.FoodCell)
			case g.State.Board[x][y] == 999:
				builder.WriteString(util.RED + g.State.Config.MazeChar + util.BLACK)
			case g.State.Board[x][y] < -1:
				switch util.PowerUpType(g.State.Board[x][y]) {
				case util.SpeedUp:
					builder.WriteString("⚡")
				case util.SlowDown:
					builder.WriteString("⏳")
				case util.GhostMode:
					builder.WriteString("👻")
				case util.ExtraLength:
					builder.WriteString("🔄")
				case util.DoublePoints:
					builder.WriteString("💎")
				}
			case g.State.Board[x][y] > 0:
				for _, powerup := range g.PowerMgr.ActivePowerUps {
					switch powerup.Type {
					case util.SpeedUp:
						builder.WriteString(util.YELLOW)
					case util.SlowDown:
						builder.WriteString(util.BLUE)
					}
				}
				if g.State.Board[x][y] == 1 {
					builder.WriteString(g.State.Config.SnakeHead)
				} else {
					builder.WriteString(g.State.Config.SnakeCell)
				}
				builder.WriteString(util.BLACK)
			}
		}
		r.decideColor(builder, g.PowerMgr.GhostMode)
		builder.WriteString(g.State.Config.BorderChar + "\n")
		builder.WriteString(util.BLACK)
	}
}

func (r *Renderer) renderTopAndBottomBorder(builder *strings.Builder, isGhostMode bool) {
	r.decideColor(builder, isGhostMode)
	builder.WriteString(strings.Repeat(" ", r.Config.OffsetX-1))
	builder.WriteString(strings.Repeat(r.Config.BorderChar, 2*(r.Config.TermWidth+1)))
	builder.WriteString(util.BLACK)
	builder.WriteString("\n")
}

func (r *Renderer) renderTopOffset(builder *strings.Builder) {
	for x := 0; x < r.Config.OffsetY-2; x++ {
		builder.WriteString("\n")
	}
}

func (r *Renderer) renderScore(builder *strings.Builder, score int) {
	builder.WriteString("\n" + strings.Repeat(" ", r.Config.OffsetX-1) + "Score: " + strconv.Itoa(score) + "\n")
}

func (r *Renderer) renderActiveEffects(builder *strings.Builder, g *Game) {
	if len(g.PowerMgr.ActivePowerUps) == 0 {
		builder.WriteString(strings.Repeat(" ", r.Config.OffsetX-1))
		builder.WriteString("Active Effects: None" + strings.Repeat(" ", g.State.Config.TermWidth))
		return
	}

	builder.WriteString(strings.Repeat(" ", r.Config.OffsetX-1))
	builder.WriteString("Active Effects: ")

	for _, powerup := range g.PowerMgr.ActivePowerUps {
		remaining := time.Until(powerup.EndTime).Seconds()
		if remaining <= 0 {
			continue
		}

		symbol := ""
		effect := ""
		switch powerup.Type {
		case util.SpeedUp:
			symbol = "⚡"
			effect = "Speed Up"
		case util.SlowDown:
			symbol = "⏳"
			effect = "Slow Down"
		case util.GhostMode:
			symbol = "👻"
			effect = "Ghost Mode"
		case util.ExtraLength:
			symbol = "🔄"
			effect = "Extra Length"
		case util.DoublePoints:
			symbol = "💎"
			effect = "Double Points"
		}

		builder.WriteString(fmt.Sprintf("%s %s (%.1fs) ", symbol, effect, remaining))

	}
	builder.WriteString(strings.Repeat(" ", g.State.Config.TermWidth-2*g.State.Config.OffsetX))
	builder.WriteString("\n")
}
