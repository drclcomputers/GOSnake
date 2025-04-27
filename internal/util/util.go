package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/term"
)

func NewGameConfig() *GameConfig {
	width, height := GetBoardSize()
	offsetX, offsetY := CalculateOffsets()

	return &GameConfig{
		TermWidth:  width,
		TermHeight: height,
		OffsetX:    offsetX,
		OffsetY:    offsetY,
		Speed:      200 * time.Millisecond,
		BorderChar: "#",
		MazeChar:   "##",
		EmptyCell:  "  ",
		SnakeCell:  "()",
		SnakeHead:  ":)",
		FoodCell:   "🍎",
	}
}

func NewSnake() *Snake {
	return &Snake{
		Headx:     0,
		Heady:     0,
		Direction: 2,
		Length:    1,
	}
}

func GetBoardSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Warning: Unable to determine terminal size, using default 80x24.")
		width, height = 80, 24
	}
	width = (width / 2) * 2 / 3
	height = (height - 1) * 2 / 3
	if width < 40 {
		width = 40
	}
	if height < 20 {
		height = 20
	}
	return width, height
}

func CalculateOffsets() (int, int) {
	fullWidth, fullHeight, err := term.GetSize(int(os.Stdout.Fd()))
	termWidth, termHeight := GetBoardSize()
	if err != nil {
		return 1, 1
	}

	offsetX := (fullWidth - (termWidth * 2)) / 2
	offsetY := (fullHeight - termHeight) / 2

	return offsetX, offsetY
}

func InitializeBoard(width, height int) [][]int {
	board := make([][]int, height)
	for i := range board {
		board[i] = make([]int, width)
	}
	return board
}

func HideCursor() {
	fmt.Print("\033[?25l")
}

func ShowCursor() {
	fmt.Print("\033[?25h")
}

func GoAtTopLeft() { fmt.Print("\033[H") }

func ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Failed to clear screen: %v\n", err)
	}
}

func MoveCursorTo(x, y int) {
	fmt.Printf("\033[%d;%dH", x, y)
}
