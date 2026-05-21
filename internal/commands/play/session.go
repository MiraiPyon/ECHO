package play

import (
	"math"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

var playbackSpeeds = []float64{0.25, 0.5, 0.75, 1, 1.25, 1.5, 1.75, 2}

type playbackSession struct {
	filePath string

	streamer beep.StreamSeekCloser
	format   beep.Format
	total    time.Duration

	stateMu     sync.RWMutex
	ctrl        *beep.Ctrl
	volume      *effects.Volume
	paused      bool
	muted       bool
	volumeLevel float64
	speedIndex  int
	finished    bool
	active      bool
	done        chan struct{}
	chainID     uint64
}

type playbackSnapshot struct {
	FilePath    string
	Current     time.Duration
	Total       time.Duration
	Paused      bool
	Muted       bool
	VolumeLevel float64
	Speed       float64
	Finished    bool
	Active      bool
}

func newPlaybackSession(filePath string) (*playbackSession, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	streamer, format, err := mp3.Decode(file)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		_ = streamer.Close()
		_ = file.Close()
		return nil, err
	}

	return &playbackSession{
		filePath:    filePath,
		streamer:    streamer,
		format:      format,
		total:       format.SampleRate.D(streamer.Len()),
		volumeLevel: 0,
		speedIndex:  nearestSpeedIndex(1),
	}, nil
}

func (s *playbackSession) startPlayback() {
	s.stateMu.Lock()
	done := make(chan struct{})
	s.done = done
	s.active = true
	s.finished = false
	s.stateMu.Unlock()

	s.playChain(done)
}

func (s *playbackSession) playChain(done chan struct{}) {
	s.stateMu.Lock()
	chainID := s.chainID + 1
	s.chainID = chainID
	paused := s.paused
	speed := playbackSpeeds[speedIndexSafe(s.speedIndex)]
	ctrl := &beep.Ctrl{Streamer: s.streamer, Paused: paused}
	resampler := beep.ResampleRatio(4, speed, ctrl)
	vol := &effects.Volume{
		Streamer: resampler,
		Base:     2,
		Volume:   s.volumeLevel,
		Silent:   s.muted,
	}
	s.ctrl = ctrl
	s.volume = vol
	s.stateMu.Unlock()

	speaker.Clear()

	speaker.Play(beep.Seq(vol, beep.Callback(func() {
		s.stateMu.Lock()
		if s.chainID != chainID {
			s.stateMu.Unlock()
			return
		}
		s.active = false
		s.finished = true
		s.paused = true
		s.stateMu.Unlock()
		close(done)
	})))
}

func (s *playbackSession) closePlayback() {
	s.stateMu.Lock()
	s.chainID++
	ctrl := s.ctrl
	streamer := s.streamer
	s.ctrl = nil
	s.volume = nil
	s.done = nil
	s.stateMu.Unlock()

	speaker.Lock()
	if ctrl != nil {
		ctrl.Streamer = nil
	}
	speaker.Unlock()
	speaker.Clear()

	if streamer != nil {
		_ = streamer.Close()
	}
	speaker.Close()
}

func (s *playbackSession) Close() {
	s.closePlayback()
}

func (s *playbackSession) snapshot() playbackSnapshot {
	s.stateMu.RLock()
	streamer := s.streamer
	paused := s.paused
	muted := s.muted
	volumeLevel := s.volumeLevel
	speedIndex := s.speedIndex
	finished := s.finished
	active := s.active
	total := s.total
	filePath := s.filePath
	format := s.format
	s.stateMu.RUnlock()

	current := time.Duration(0)
	if streamer != nil {
		speaker.Lock()
		current = format.SampleRate.D(streamer.Position())
		speaker.Unlock()
	}

	if current < 0 {
		current = 0
	}
	if current > total {
		current = total
	}

	return playbackSnapshot{
		FilePath:    filePath,
		Current:     current,
		Total:       total,
		Paused:      paused,
		Muted:       muted,
		VolumeLevel: volumeLevel,
		Speed:       playbackSpeeds[speedIndexSafe(speedIndex)],
		Finished:    finished,
		Active:      active,
	}
}

func (s *playbackSession) togglePause() {
	s.stateMu.Lock()
	active := s.active
	if s.finished {
		s.stateMu.Unlock()
		if s.currentPosition() >= s.total {
			_ = s.seekTo(0)
		}
		s.stateMu.Lock()
		s.paused = false
		s.finished = false
		s.stateMu.Unlock()
		s.startPlayback()
		return
	}

	if !active {
		s.paused = false
		s.stateMu.Unlock()
		s.startPlayback()
		return
	}

	s.paused = !s.paused
	paused := s.paused
	ctrl := s.ctrl
	s.stateMu.Unlock()

	speaker.Lock()
	if ctrl != nil {
		ctrl.Paused = paused
	}
	speaker.Unlock()
}

func (s *playbackSession) toggleMute() {
	s.stateMu.Lock()
	s.muted = !s.muted
	muted := s.muted
	vol := s.volume
	s.stateMu.Unlock()

	speaker.Lock()
	if vol != nil {
		vol.Silent = muted
	}
	speaker.Unlock()
}

func (s *playbackSession) setSpeed(speed float64) {
	index := nearestSpeedIndex(speed)

	s.stateMu.Lock()
	if index == s.speedIndex {
		s.stateMu.Unlock()
		return
	}

	s.speedIndex = index
	active := s.active
	done := s.done
	s.stateMu.Unlock()

	if active && done != nil {
		s.playChain(done)
	}
}

func (s *playbackSession) adjustVolume(delta float64) {
	s.stateMu.Lock()
	s.volumeLevel = clampFloat(s.volumeLevel+delta, -2, 2)
	level := s.volumeLevel
	vol := s.volume
	s.stateMu.Unlock()

	speaker.Lock()
	if vol != nil {
		vol.Volume = level
		vol.Silent = false
	}
	speaker.Unlock()
}

func (s *playbackSession) seekBy(delta time.Duration) error {
	return s.seekTo(s.currentPosition() + delta)
}

func (s *playbackSession) seekTo(target time.Duration) error {
	total := s.total
	if target < 0 {
		target = 0
	}
	if target > total {
		target = total
	}

	s.stateMu.RLock()
	active := s.active
	done := s.done
	s.stateMu.RUnlock()

	speaker.Lock()
	err := s.streamer.Seek(s.format.SampleRate.N(target))
	speaker.Unlock()
	if err != nil {
		return err
	}

	s.stateMu.Lock()
	if s.finished {
		s.finished = false
	}
	s.stateMu.Unlock()

	if active && done != nil {
		s.playChain(done)
	}

	return nil
}

func (s *playbackSession) currentPosition() time.Duration {
	s.stateMu.RLock()
	streamer := s.streamer
	format := s.format
	s.stateMu.RUnlock()

	if streamer == nil {
		return 0
	}

	speaker.Lock()
	pos := format.SampleRate.D(streamer.Position())
	speaker.Unlock()
	return pos
}

func (s *playbackSession) doneChannel() <-chan struct{} {
	s.stateMu.RLock()
	defer s.stateMu.RUnlock()
	return s.done
}

func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func playSimple(filePath string) error {
	session, err := newPlaybackSession(filePath)
	if err != nil {
		return err
	}
	defer session.Close()

	session.startPlayback()
	<-session.doneChannel()
	return nil
}

func speedIndexSafe(index int) int {
	if index < 0 {
		return 0
	}
	if index >= len(playbackSpeeds) {
		return len(playbackSpeeds) - 1
	}
	return index
}

func nearestSpeedIndex(speed float64) int {
	bestIndex := 0
	bestDelta := math.Abs(playbackSpeeds[0] - speed)
	for i := 1; i < len(playbackSpeeds); i++ {
		delta := math.Abs(playbackSpeeds[i] - speed)
		if delta < bestDelta {
			bestDelta = delta
			bestIndex = i
		}
	}
	return bestIndex
}

// AvailablePlaybackSpeeds returns a copy of the supported playback speed presets.
func AvailablePlaybackSpeeds() []float64 {
	speeds := make([]float64, len(playbackSpeeds))
	copy(speeds, playbackSpeeds)
	return speeds
}

// NearestSpeedIndex returns the preset index closest to the requested speed.
func NearestSpeedIndex(speed float64) int {
	return nearestSpeedIndex(speed)
}

// SpeedIndexSafe clamps an index so it always points to a valid preset.
func SpeedIndexSafe(index int) int {
	return speedIndexSafe(index)
}

// FormatSpeedLabel formats a playback speed for display.
func FormatSpeedLabel(speed float64) string {
	return formatSpeedLabel(speed)
}
