package game

import (
	"fmt"
	"gosnake/internal/util"
	"os"

	"github.com/eiannone/keyboard"
)

func (g *Game) SetRelaxedMode(relaxed bool) {
	g.State.RelaxedMode = relaxed
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
	switch g.State.Config.Mode {
	case util.Normal:
		fmt.Println("Normal - Classic Snake gameplay with increasing speed")
	case util.NoWalls:
		fmt.Println("No Walls - Snake can pass through borders")
	case util.Maze:
		fmt.Println("Maze - Navigate through randomly generated obstacles")
	case util.PowerUps:
		fmt.Println("Power-ups - Collect special items for unique abilities:")
		fmt.Println("  ⚡ Speed Up   ⏳ Slow Down   👻 Ghost Mode")
		fmt.Println("  🔄 Extra Length   💎 Double Points")
	}
	if g.State.RelaxedMode {
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
	switch g.State.Config.Mode {
	case util.Maze:
		g.generateMaze()
	case util.PowerUps:
		g.spawnPowerUp()
	}

	g.placeFood()
	g.State.Board[g.State.Snake.Headx][g.State.Snake.Heady] = 1
	go g.pollInput()
}
