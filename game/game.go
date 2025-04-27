// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package game

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gosnake/internal/util"

	"github.com/eiannone/keyboard"
)

type Game struct {
	config          *util.GameConfig
	snake           *util.Snake
	board           [][]int
	score           int
	inputChan       chan keyboard.KeyEvent
	exitGame        bool
	exitCode        int
	pauseGame       bool
	relaxedMode     bool
	ghostMode       bool
	pointMultiplier int
	activePowerUps  []*PowerUp
	sound           *SoundManager
}

func NewGame(config *util.GameConfig) *Game {
	return &Game{
		config:          config,
		snake:           util.NewSnake(),
		board:           util.InitializeBoard(config.TermWidth, config.TermHeight),
		score:           0,
		inputChan:       make(chan keyboard.KeyEvent, 10),
		exitGame:        false,
		exitCode:        0,
		pauseGame:       false,
		relaxedMode:     false,
		ghostMode:       false,
		pointMultiplier: 1,
		activePowerUps:  make([]*PowerUp, 0),
		sound:           NewSoundManager(false),
	}
}

func killSig() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		util.ShowCursor()
		keyboard.Close()
		os.Exit(0)
	}()
}

func (g *Game) SetRelaxedMode(relaxed bool) {
	g.relaxedMode = relaxed
}

func (g *Game) InitSound(enableSound bool) {
	g.sound = NewSoundManager(enableSound)

	if enableSound {
		go PlayMusic("./assets/background.mp3", -1) // Loop forever
	}
}

func (g *Game) Welcome() {
	util.ClearScreen()

	fmt.Println("Welcome to GOSNAKE!")
	fmt.Println()
	fmt.Println("Controls: WASD or Arrow to change direction, 'p' to pause or 'q' to quit.")
	fmt.Println()

	fmt.Print("Game Mode: ")
	switch g.config.Mode {
	case util.Normal:
		fmt.Println("Normal - Classic snake gameplay with increasing speed")
	case util.NoWalls:
		fmt.Println("No Walls - Snake can pass through borders")
	case util.Maze:
		fmt.Println("Maze - Navigate through randomly generated obstacles")
	case util.PowerUps:
		fmt.Println("Power-ups - Collect special items for unique abilities:")
		fmt.Println("  ⚡ Speed Up   ⏳ Slow Down   👻 Ghost Mode")
		fmt.Println("  🔄 Extra Length   💎 Double Points")
	}
	if g.relaxedMode {
		fmt.Println("Relaxed Mode: ON - Speed remains constant")
	}
	fmt.Println()

	printHighScores()
	fmt.Println()
	fmt.Println("Press any key to start...")

	var in string
	fmt.Scanln(&in)

	if in == "q" {
		os.Exit(0)
	}
}

func (g *Game) Start() {
	g.Welcome()

	killSig()

	util.ClearScreen()
	util.HideCursor()
	defer util.ShowCursor()

	if err := keyboard.Open(); err != nil {
		fmt.Println("Error initializing keyboard input:", err)
		return
	}
	defer keyboard.Close()

	g.initializeGame()
	g.runGameLoop()
}

func (g *Game) initializeGame() {
	switch g.config.Mode {
	case util.Maze:
		g.generateMaze()
	case util.PowerUps:
		g.spawnPowerUp()
	}

	g.placeFood()
	g.board[g.snake.Headx][g.snake.Heady] = 1
	go g.pollInput()
}

func (g *Game) placeFood() {
	g.clearEatenApples()

	x, y := g.getRandomEmptyPosition()
	g.board[x][y] = -1
}

func (g *Game) clearEatenApples() {
	for x := 0; x < g.config.TermHeight; x++ {
		for y := 0; y < g.config.TermWidth; y++ {
			if g.board[x][y] == -1 {
				g.board[x][y] = 0
			}
		}
	}
}

func (g *Game) detectPause() {
	for g.pauseGame {
		select {
		case event := <-g.inputChan:
			g.handleInput(event)
		}
	}
}

func (g *Game) runGameLoop() {
	ticker := time.NewTicker(g.config.Speed)
	defer ticker.Stop()

	renderer := NewRenderer(g.config)
	renderer.Render(g)

	for !g.exitGame {
		g.detectPause()

		select {
		case event := <-g.inputChan:
			g.handleInput(event)
			renderer.Render(g)
		case <-ticker.C:
			g.update()
			if !g.exitGame {
				renderer.Render(g)
			}
		}
	}

	go g.sound.PlayGameOver()

	fmt.Println("\nPress any key to continue...")

	g.writeHighScores()
}

func (g *Game) update() {
	g.moveSnake()
	col := g.checkCollision()
	if col != 0 {
		g.exitCode = col
		g.exitGame = true
		StopMusic()
		if col == 1 {
			fmt.Println("\n" + strings.Repeat(" ", g.config.OffsetX-1) + "Game Over! Snake hit a wall!")
		} else if col == 2 {
			fmt.Println("\n" + strings.Repeat(" ", g.config.OffsetX-1) + "Game Over! Snake collided with itself!")
		}
		return
	}
	g.updateBoard()
	if g.score%15 == 0 && g.score > 0 && !g.relaxedMode {
		g.config.Speed -= 15
	}
}

func (g *Game) pollInput() {
	for !g.exitGame {
		char, key, err := keyboard.GetKey()
		if err != nil {
			if g.exitGame {
				return
			}
			continue
		}
		if !g.exitGame {
			g.inputChan <- keyboard.KeyEvent{Key: key, Rune: char}
		}
	}
}

func (g *Game) handleInput(event keyboard.KeyEvent) {
	switch {
	case event.Key == keyboard.KeyEsc || event.Rune == 'q' || event.Rune == 'Q':
		g.exitGame = true
	case (event.Key == keyboard.KeyArrowLeft || event.Rune == 'a') && g.snake.Direction != 2 && !g.pauseGame:
		g.snake.Direction = 4
		g.config.SnakeHead = "(:"
	case (event.Key == keyboard.KeyArrowRight || event.Rune == 'd') && g.snake.Direction != 4 && !g.pauseGame:
		g.snake.Direction = 2
		g.config.SnakeHead = ":)"
	case (event.Key == keyboard.KeyArrowDown || event.Rune == 's') && g.snake.Direction != 1 && !g.pauseGame:
		g.snake.Direction = 3
		g.config.SnakeHead = "()"
	case (event.Key == keyboard.KeyArrowUp || event.Rune == 'w') && g.snake.Direction != 3 && !g.pauseGame:
		g.snake.Direction = 1
		g.config.SnakeHead = "()"
	case event.Rune == 'p' || event.Rune == 'P':
		g.pauseGame = !g.pauseGame
	}
}

func readHighScores() []int {
	file, _ := os.ReadFile("score.txt")
	lines := strings.Split(string(file), "\n")
	scores := []int{}
	for _, line := range lines {
		if score, err := strconv.Atoi(line); err == nil {
			scores = append(scores, score)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(scores)))
	return scores
}

func printHighScores() {
	fmt.Println("Top Scores:")
	scores := readHighScores()
	for i, score := range scores {
		fmt.Printf("%d. %d\n", i+1, score)
		if i == 4 {
			break
		}
	}
}

func (g *Game) writeHighScores() {
	if g.score == 0 {
		return
	}

	file, err := os.OpenFile("score.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error saving high score:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%d\n", g.score))
	if err != nil {
		fmt.Println("Error writing high score:", err)
	}
}
