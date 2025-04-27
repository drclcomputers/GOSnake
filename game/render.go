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
	config *util.GameConfig
}

func NewRenderer(config *util.GameConfig) *Renderer {
	return &Renderer{config: config}
}

func (r *Renderer) Render(g *Game) {
	var builder strings.Builder

	r.renderTopOffset(&builder)
	r.renderTopAndBottomBorder(&builder)
	r.renderBoard(&builder, g)
	r.renderTopAndBottomBorder(&builder)
	r.renderScore(&builder, g.score)
	if g.config.Mode == util.PowerUps {
		r.renderActiveEffects(&builder, g)
	}

	util.GoAtTopLeft()
	fmt.Print(builder.String())
}

func (r *Renderer) decideColor(builder *strings.Builder) {
	switch r.config.Mode {
	case util.NoWalls:
		builder.WriteString(util.GREEN)
	case util.Maze, util.Normal:
		builder.WriteString(util.RED)
	default:
		builder.WriteString(util.BLACK)
	}
}

func (r *Renderer) renderBoard(builder *strings.Builder, g *Game) {
	for x := 0; x < g.config.TermHeight; x++ {
		r.decideColor(builder)
		builder.WriteString(strings.Repeat(" ", g.config.OffsetX-1) + g.config.BorderChar)
		builder.WriteString(util.BLACK)

		for y := 0; y < g.config.TermWidth; y++ {
			switch {
			case g.board[x][y] == 0:
				builder.WriteString(g.config.EmptyCell)
			case g.board[x][y] == -1:
				builder.WriteString(g.config.FoodCell)
			case g.board[x][y] == 999:
				builder.WriteString(util.RED + g.config.MazeChar + util.BLACK)
			case g.board[x][y] < -1:
				switch PowerUpType(g.board[x][y]) {
				case SpeedUp:
					builder.WriteString("⚡ ")
				case SlowDown:
					builder.WriteString("⏳ ")
				case GhostMode:
					builder.WriteString("👻 ")
				case ExtraLength:
					builder.WriteString("🔄 ")
				case DoublePoints:
					builder.WriteString("💎 ")
				}
			case g.board[x][y] == 1:
				builder.WriteString(g.config.SnakeHead)
			default:
				builder.WriteString(g.config.SnakeCell)
			}
		}
		r.decideColor(builder)
		builder.WriteString(g.config.BorderChar + "\n")
		builder.WriteString(util.BLACK)
	}
}

func (r *Renderer) renderTopAndBottomBorder(builder *strings.Builder) {
	r.decideColor(builder)
	builder.WriteString(strings.Repeat(" ", r.config.OffsetX-1))
	builder.WriteString(strings.Repeat(r.config.BorderChar, 2*(r.config.TermWidth+1)))
	builder.WriteString(util.BLACK)
	builder.WriteString("\n")
}

func (r *Renderer) renderTopOffset(builder *strings.Builder) {
	for x := 0; x < r.config.OffsetY-2; x++ {
		builder.WriteString("\n")
	}
}

func (r *Renderer) renderScore(builder *strings.Builder, score int) {
	builder.WriteString("\n" + strings.Repeat(" ", r.config.OffsetX-1) + "Score: " + strconv.Itoa(score) + "\n")
}

func (r *Renderer) renderActiveEffects(builder *strings.Builder, g *Game) {
	if len(g.activePowerUps) == 0 {
		builder.WriteString(strings.Repeat(" ", r.config.OffsetX-1))
		builder.WriteString("Active Effects: None" + strings.Repeat(" ", g.config.TermWidth))
		return
	}

	builder.WriteString(strings.Repeat(" ", r.config.OffsetX-1))
	builder.WriteString("Active Effects: ")

	for _, powerup := range g.activePowerUps {
		remaining := time.Until(powerup.EndTime).Seconds()
		if remaining <= 0 {
			continue
		}

		symbol := ""
		effect := ""
		switch powerup.Type {
		case SpeedUp:
			symbol = "⚡"
			effect = "Speed Up"
		case SlowDown:
			symbol = "⏳"
			effect = "Slow Down"
		case GhostMode:
			symbol = "👻"
			effect = "Ghost Mode"
		case ExtraLength:
			symbol = "🔄"
			effect = "Extra Length"
		case DoublePoints:
			symbol = "💎"
			effect = "Double Points"
		}

		builder.WriteString(fmt.Sprintf("%s %s (%.1fs) ", symbol, effect, remaining))

	}
	builder.WriteString(strings.Repeat(" ", g.config.TermWidth-2*g.config.OffsetX))
	builder.WriteString("\n")
}
