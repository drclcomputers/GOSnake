// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package game

import (
	"fmt"
	"gosnake/internal/util"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
)

func (g *Game) pollInput() {
	for !g.State.ExitGame {
		char, key, err := keyboard.GetKey()
		if err != nil {
			if g.State.ExitGame {
				return
			}
			continue
		}
		if !g.State.ExitGame {
			g.inputChan <- keyboard.KeyEvent{Key: key, Rune: char}
		}
	}
}

func (g *Game) handleInput(event keyboard.KeyEvent) {
	switch {
	case event.Key == keyboard.KeyEsc || event.Rune == 'q' || event.Rune == 'Q':
		g.State.ExitGame = true
	case (event.Key == keyboard.KeyArrowLeft || event.Rune == 'a') && g.State.Snake.Direction != 2 && !g.State.PauseGame:
		g.State.Snake.Direction = util.DirectionLeft
		g.State.Config.SnakeHead = "(:"
	case (event.Key == keyboard.KeyArrowRight || event.Rune == 'd') && g.State.Snake.Direction != 4 && !g.State.PauseGame:
		g.State.Snake.Direction = util.DirectionRight
		g.State.Config.SnakeHead = ":)"
	case (event.Key == keyboard.KeyArrowDown || event.Rune == 's') && g.State.Snake.Direction != 1 && !g.State.PauseGame:
		g.State.Snake.Direction = util.DirectionDown
		g.State.Config.SnakeHead = "()"
	case (event.Key == keyboard.KeyArrowUp || event.Rune == 'w') && g.State.Snake.Direction != 3 && !g.State.PauseGame:
		g.State.Snake.Direction = util.DirectionUp
		g.State.Config.SnakeHead = "()"
	case event.Rune == 'p' || event.Rune == 'P':
		g.State.PauseGame = !g.State.PauseGame
	}
}

func readHighScores() []int {
	file, _ := os.ReadFile("Score.txt")
	lines := strings.Split(string(file), "\n")
	scores := []int{}
	for _, line := range lines {
		if Score, err := strconv.Atoi(line); err == nil {
			scores = append(scores, Score)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(scores)))
	return scores
}

func printHighScores() {
	fmt.Println("Top Scores:")
	scores := readHighScores()
	for i, Score := range scores {
		fmt.Printf("%d. %d\n", i+1, Score)
		if i == 4 {
			break
		}
	}
}

func (g *Game) writeHighScores() {
	if g.State.Score == 0 {
		return
	}

	file, err := os.OpenFile("Score.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error saving high Score:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d\n", g.State.Score))
	if err != nil {
		fmt.Println("Error writing high Score:", err)
	}
}
