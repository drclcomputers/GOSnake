// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package util

import "time"

type Snake struct {
	Headx     int
	Heady     int
	Direction int // 1 - up, 2 - right, 3 - down, 4 - left
	Length    int
}

type GameMode int

const (
	Normal   GameMode = iota
	NoWalls           // Snake passes through walls
	Maze              // Has obstacles
	PowerUps          // Includes power-ups
)

type GameConfig struct {
	TermWidth  int
	TermHeight int
	OffsetX    int
	OffsetY    int
	Speed      time.Duration
	BorderChar string
	MazeChar   string
	EmptyCell  string
	SnakeCell  string
	SnakeHead  string
	FoodCell   string
	PowerUp    string
	Mode       GameMode
	Obstacles  []Position
}

type Position struct {
	X, Y int
}

const (
	BLACK   = "\033[0m"
	RED     = "\033[31m"
	GREEN   = "\033[32m"
	YELLOW  = "\033[33m"
	BLUE    = "\033[34m"
	MAGENTA = "\033[35m"
	CYAN    = "\033[36m"
	WHITE   = "\033[37m"
)

const VER = "v0.6"
