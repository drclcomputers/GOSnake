package cmd

import (
	"gosnake/game"
	"gosnake/internal/util"
	"time"

	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Start playing Snake",
	Run: func(cmd *cobra.Command, args []string) {
		config := util.NewGameConfig()
		config.Speed = time.Duration(speed) * time.Millisecond

		switch gameMode {
		case "nowalls":
			config.Mode = util.NoWalls
		case "maze":
			config.Mode = util.Maze
		case "powerups":
			config.Mode = util.PowerUps
		default:
			config.Mode = util.Normal
		}

		game := game.NewGame(config)
		game.SetRelaxedMode(relaxed)
		game.InitSound(!noSound)
		game.Start()
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}
