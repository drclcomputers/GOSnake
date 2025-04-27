# 🐍 GoSnake

A feature-rich terminal-based Snake game written in Go, featuring multiple game modes, power-ups, and sound effects!

## ✨ Features

### 🎮 Game Modes
- **Normal Mode**: Classic snake gameplay with increasing speed
- **No Walls Mode**: Snake can pass through borders
- **Maze Mode**: Navigate through randomly generated obstacles
- **Power-ups Mode**: Collect special items for unique abilities

### 🎲 Power-ups
- ⚡ **Speed Up**: Temporarily increase snake's speed
- ⏳ **Slow Down**: Temporarily decrease snake's speed
- 👻 **Ghost Mode**: Pass through walls and obstacles
- 🔄 **Extra Length**: Instantly grow longer
- 💎 **Double Points**: Score multiplier

### 🎯 Additional Features
- High score tracking
- Sound effects and background music
- Pause functionality
- Relaxed mode option (constant speed)
- Colorful terminal UI
- Custom starting speed

## 🚀 Installation

1. Make sure you have Go 1.24 or later installed
2. Clone the repository:
```bash
git clone https://github.com/yourusername/go-snake.git
cd go-snake
```

3. Build and run:
```bash
go build
./gosnake play
```

## 🎮 Controls

- **Arrow Keys** or **WASD**: Control snake direction
- **P**: Pause game
- **Q** or **ESC**: Quit game

## 🎯 Command Line Options

```bash
# Start normal game
./gosnake play

# Play with no walls
./gosnake play --mode nowalls

# Play maze mode
./gosnake play --mode maze

# Play with power-ups
./gosnake play --mode powerups

# Play in relaxed mode (constant speed)
./gosnake play -r ( or --relaxed)

# Set custom initial speed (milliseconds)
./gosnake play --speed 150

# Disable sound (music and sound effects)
./gosnake play --no-sound
```

## 🏆 Scoring

- Each food item: 1 point
- With Double Points power-up: 2 points per food
- High scores are automatically saved
- View top 5 scores at game start

## 🎨 Terminal Display

```
##################
#                #
#  :)🍎          #
#  ()            #
#      ⚡        #
#                #
##################
Score: 5
Active Effects: ⚡ Speed Up (5.2s)
```

## 🛠 Technical Requirements

- Go 1.24.2 or later
- Terminal with ANSI color support
- Audio output capability (for sound effects)

## 🐛 Known Bugs
- The snake can sometimes become schizophrenic and see two apples. This effect may or may not disappear after eating one of them. 🍏

## 📝 License

MIT License - feel free to use and modify!

## 🤝 Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

## 🙏 Acknowledgments

- Built with Go and various awesome Go packages
- Inspired by the classic Snake game
- Special thanks to the Go community
- Music by [Nicholas Panek](https://pixabay.com/users/nickpanek620-38266323/?utm_source=link-attribution&utm_medium=referral&utm_campaign=music&utm_content=318059) from [Pixabay](https://pixabay.com/music//?utm_source=link-attribution&utm_medium=referral&utm_campaign=music&utm_content=318059)

---

Made with ❤️ by @drclcomputers
