// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package cmd

import (
	"fmt"
	"os"

	"gosnake/internal/util"

	"github.com/spf13/cobra"
)

var (
	gameMode string
	speed    int
	relaxed  bool
	noSound  bool
	rootCmd  = &cobra.Command{
		Use:   "gosnake",
		Short: "A terminal-based Snake game written in Go",
		Long: `A terminal-based Snake game with various options:
- Normal mode: Speed increases as you collect food
- Maze mode: Navigate through a randomly generated maze, collectig food
- No Walls Mode: There are no borders
- PowerUps: Enhance your abilities with powerups
- Relaxed mode: Speed remains constant accross all game modes
- Custom starting speed`,
		Version: util.VER,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&gameMode, "mode", "m", "normal", "Game mode (normal, nowalls, maze, powerups)")
	rootCmd.PersistentFlags().IntVarP(&speed, "speed", "s", 200, "Initial game speed (milliseconds)")
	rootCmd.PersistentFlags().BoolVarP(&relaxed, "relaxed", "r", false, "Enable relaxed mode (constant speed)")
	rootCmd.PersistentFlags().BoolVar(&noSound, "no-sound", false, "Disable sound")
}
