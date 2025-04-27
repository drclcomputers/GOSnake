// Copyright (c) 2025 @drclcomputers. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package game

import (
	"fmt"
	"gosnake/internal/util"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gen2brain/beeep"
)

type SoundManager struct {
	enabled bool
}

var musicCtrl *beep.Ctrl
var musicDone chan bool

func NewSoundManager(enableSound bool) *SoundManager {
	return &SoundManager{
		enabled: enableSound,
	}
}

func (s *SoundManager) PlayFoodEaten() {
	if s == nil || !s.enabled {
		return
	}
	go beeep.Beep(880, 200)
}

func (s *SoundManager) PlayPowerUpCollected() {
	if s == nil || !s.enabled {
		return
	}
	go func() {
		beeep.Beep(587.33, 200)
		time.Sleep(200 * time.Millisecond)
		beeep.Beep(880, 200)
	}()
}

func (s *SoundManager) PlayGameOver() {
	if s == nil || !s.enabled {
		return
	}
	go beeep.Beep(320, 300)
}

func (s *SoundManager) ToggleSound() {
	s.enabled = !s.enabled
}

func PlayMusic(sound string, times int) {
	if !util.CheckSpeaker() {
		fmt.Println("No speaker found. Sound disabled!")
		return
	}

	f, err := os.Open(sound)
	if err != nil {
		return
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return
	}

	musicDone = make(chan bool)

	loop := beep.Loop(times, streamer)
	musicCtrl = &beep.Ctrl{Streamer: loop, Paused: false}

	speaker.Play(beep.Seq(
		musicCtrl,
		beep.Callback(func() {
			musicDone <- true
		}),
	))

	<-musicDone
}

func StopMusic() {
	if musicCtrl != nil {
		speaker.Lock()
		musicCtrl.Paused = true
		speaker.Unlock()

		if musicDone != nil {
			close(musicDone)
		}
	}
}
