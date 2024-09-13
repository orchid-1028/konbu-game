package main

import (
	"bytes"
	_ "embed"
	_ "image/png"
	"io"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

// Player represents the current audio state.
type Player struct {
	game         *Game
	audioContext *audio.Context
	audioPlayer  *audio.Player
	current      time.Duration
	total        time.Duration
	seBytes      []byte
	seCh         chan []byte
	volume128    int
}

const (
	sampleRate = 48000
)

//go:embed assets/master_5minutes.mp3
var BGM_mp3 []byte

func NewPlayer(game *Game, audioContext *audio.Context) (*Player, error) {
	type audioStream interface {
		io.ReadSeeker
		Length() int64
	}

	const bytesPerSample = 4 // TODO: This should be defined in audio package

	var s audioStream

	var err error
	s, err = mp3.DecodeWithoutResampling(bytes.NewReader(BGM_mp3))
	if err != nil {
		return nil, err
	}
	p, err := audioContext.NewPlayer(s)
	if err != nil {
		return nil, err
	}
	player := &Player{
		game:         game,
		audioContext: audioContext,
		audioPlayer:  p,
		total:        time.Second * time.Duration(s.Length()) / bytesPerSample / sampleRate,
		volume128:    128,
		seCh:         make(chan []byte),
	}
	if player.total == 0 {
		player.total = 1
	}

	player.audioPlayer.Play()
	return player, nil
}

func NewGame() (*Game, error) {
	audioContext := audio.NewContext(sampleRate)

	g := &Game{
		musicPlayerCh: make(chan *Player),
		errCh:         make(chan error),
	}

	m, err := NewPlayer(g, audioContext)
	if err != nil {
		return nil, err
	}

	g.musicPlayer = m
	return g, nil
}
