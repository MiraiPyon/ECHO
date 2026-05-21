package play

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ciferia/echo/internal/tui"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type playerView struct {
	*tview.Box
	session               *playbackSession
	app                   *tview.Application
	requestBackToMainMenu func()
	showSpeedSelector     func()
}

type playerLayout struct {
	cardX, cardY, cardW, cardH int
	controlsY                  int
	playX, timeX, barX         int
	barW                       int
	volumeX, speedX, backX     int
	timeW                      int
}

func launchPlayerUI(filePath string) error {
	session, err := newPlaybackSession(filePath)
	if err != nil {
		return err
	}
	defer session.Close()

	app := tview.NewApplication().EnableMouse(true)
	pages := tview.NewPages()
	var exitErr error

	player := &playerView{
		Box:     tview.NewBox().SetBackgroundColor(tcell.NewRGBColor(9, 11, 17)),
		session: session,
		app:     app,
	}
	player.Box.SetBorder(false)

	player.requestBackToMainMenu = func() {
		exitErr = tui.ErrBackToMainMenu
		app.Stop()
	}

	player.showSpeedSelector = func() {
		if pages.HasPage("speed") {
			return
		}

		snap := session.snapshot()
		currentIndex := nearestSpeedIndex(snap.Speed)

		list := tview.NewList()
		list.ShowSecondaryText(false)
		list.SetBorder(true)
		list.SetTitle(" Speed ")
		list.SetHighlightFullLine(true)
		list.SetWrapAround(true)
		list.SetMainTextColor(tcell.ColorWhite)
		list.SetShortcutColor(tcell.ColorYellow)
		list.SetSelectedBackgroundColor(tcell.NewRGBColor(34, 197, 94))
		list.SetSelectedTextColor(tcell.ColorBlack)

		for i, speed := range playbackSpeeds {
			list.AddItem(formatSpeedLabel(speed), "", rune('1'+i), nil)
		}
		list.SetCurrentItem(currentIndex)

		closeOverlay := func() {
			pages.RemovePage("speed")
			app.SetFocus(player)
		}

		list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
			if index >= 0 && index < len(playbackSpeeds) {
				session.setSpeed(playbackSpeeds[index])
			}
			closeOverlay()
		})
		list.SetDoneFunc(closeOverlay)

		panel := tview.NewFlex().SetDirection(tview.FlexColumn)
		panel.AddItem(nil, 0, 1, false)
		panel.AddItem(list, 18, 1, true)
		panel.AddItem(nil, 0, 1, false)

		overlay := tview.NewFlex().SetDirection(tview.FlexRow)
		overlay.AddItem(nil, 0, 1, false)
		overlay.AddItem(panel, 12, 1, true)
		overlay.AddItem(nil, 0, 1, false)

		pages.AddPage("speed", overlay, true, true)
		app.SetFocus(list)
	}

	pages.AddPage("player", player, true, true)
	app.SetRoot(pages, true)
	app.SetFocus(player)
	session.startPlayback()

	stopTicker := make(chan struct{})
	go func() {
		ticker := time.NewTicker(120 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				app.QueueUpdateDraw(func() {})
			case <-stopTicker:
				return
			}
		}
	}()
	defer close(stopTicker)

	if err := app.Run(); err != nil {
		return err
	}

	return exitErr
}

func (p *playerView) Draw(screen tcell.Screen) {
	p.Box.DrawForSubclass(screen, p)

	snap := p.session.snapshot()
	layout := p.computeLayout(snap)
	cardBg := tcell.NewRGBColor(236, 239, 245)
	cardBorder := tcell.NewRGBColor(194, 202, 214)
	primary := tcell.NewRGBColor(25, 31, 43)
	accent := tcell.NewRGBColor(34, 197, 94)
	track := tcell.NewRGBColor(180, 188, 200)

	drawRoundedPanel(screen, layout.cardX, layout.cardY, layout.cardW, layout.cardH, cardBg, cardBorder)

	playIcon := "▶"
	playColor := primary
	if snap.Active && !snap.Paused {
		playIcon = "⏸"
		playColor = accent
	} else if snap.Finished {
		playIcon = "↻"
		playColor = accent
	}
	tview.Print(screen, playIcon, layout.playX, layout.controlsY, 2, tview.AlignLeft, playColor)

	timeLabel := fmt.Sprintf("%s / %s", formatDuration(snap.Current), formatDuration(snap.Total))
	tview.Print(screen, timeLabel, layout.timeX, layout.controlsY, layout.timeW, tview.AlignLeft, primary)

	progress := 0.0
	if snap.Total > 0 {
		progress = math.Max(0, math.Min(1, float64(snap.Current)/float64(snap.Total)))
	}
	drawProgressBar(screen, layout.barX, layout.controlsY, layout.barW, progress, accent, track)

	volumeIcon := "🔊"
	if snap.Muted {
		volumeIcon = "🔇"
	} else if snap.VolumeLevel < -0.5 {
		volumeIcon = "🔈"
	}
	tview.Print(screen, volumeIcon, layout.volumeX, layout.controlsY, 2, tview.AlignLeft, primary)
	tview.Print(screen, formatSpeedLabel(snap.Speed), layout.speedX, layout.controlsY, 5, tview.AlignCenter, primary)
	tview.Print(screen, "↩", layout.backX, layout.controlsY, 2, tview.AlignLeft, primary)
}

func (p *playerView) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return p.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		switch event.Key() {
		case tcell.KeyEnter:
			p.session.togglePause()
		case tcell.KeyRune:
			switch strings.ToLower(string(event.Rune())) {
			case " ":
				p.session.togglePause()
			case "b":
				if p.requestBackToMainMenu != nil {
					p.requestBackToMainMenu()
				}
			case "q":
				p.app.Stop()
			case "t":
				if p.showSpeedSelector != nil {
					p.showSpeedSelector()
				}
			case "+":
				p.session.adjustVolume(0.1)
			case "-":
				p.session.adjustVolume(-0.1)
			case "0":
				p.session.toggleMute()
			}
		case tcell.KeyLeft:
			_ = p.session.seekBy(-5 * time.Second)
		case tcell.KeyRight:
			_ = p.session.seekBy(5 * time.Second)
		case tcell.KeyUp:
			p.session.adjustVolume(0.1)
		case tcell.KeyDown:
			p.session.adjustVolume(-0.1)
		case tcell.KeyEscape:
			p.app.Stop()
		}
	})
}

func (p *playerView) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
	return p.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (bool, tview.Primitive) {
		if action != tview.MouseLeftClick {
			return false, nil
		}

		x, y := event.Position()
		if !p.InRect(x, y) {
			return false, nil
		}

		setFocus(p)
		snap := p.session.snapshot()
		layout := p.computeLayout(snap)

		if y == layout.controlsY {
			switch {
			case x >= layout.playX && x < layout.playX+2:
				p.session.togglePause()
			case x >= layout.barX && x < layout.barX+layout.barW:
				target := clickPositionToTime(x, layout.barX, layout.barW, snap.Total)
				_ = p.session.seekTo(target)
			case x >= layout.volumeX && x < layout.volumeX+2:
				p.session.toggleMute()
			case x >= layout.speedX && x < layout.speedX+5:
				if p.showSpeedSelector != nil {
					p.showSpeedSelector()
				}
			case x >= layout.backX && x < layout.backX+2:
				if p.requestBackToMainMenu != nil {
					p.requestBackToMainMenu()
				}
			}
		}

		return true, nil
	})
}

func (p *playerView) computeLayout(snap playbackSnapshot) playerLayout {
	ix, iy, iw, ih := p.GetInnerRect()
	cardH := 3
	if ih < 5 {
		cardH = maxInt(1, ih-2)
	}

	cardW := minInt(iw-4, 74)
	if cardW < 50 {
		cardW = maxInt(34, iw-2)
	}
	if cardW < 20 {
		cardW = iw
	}

	cardX := ix + (iw-cardW)/2
	cardY := iy + (ih-cardH)/2
	if cardY < iy {
		cardY = iy
	}

	controlsY := cardY + cardH/2

	timeLabel := fmt.Sprintf("%s / %s", formatDuration(snap.Current), formatDuration(snap.Total))
	timeW := tview.TaggedStringWidth(timeLabel)
	if timeW < 11 {
		timeW = 11
	}

	playX := cardX + 2
	timeX := playX + 4
	barX := timeX + timeW + 2
	backX := cardX + cardW - 3
	speedX := backX - 6
	volumeX := speedX - 3
	barW := volumeX - 2 - barX
	if barW < 8 {
		barW = 8
	}

	return playerLayout{
		cardX:     cardX,
		cardY:     cardY,
		cardW:     cardW,
		cardH:     cardH,
		controlsY: controlsY,
		playX:     playX,
		timeX:     timeX,
		barX:      barX,
		barW:      barW,
		volumeX:   volumeX,
		speedX:    speedX,
		backX:     backX,
		timeW:     timeW,
	}
}

func drawRoundedPanel(screen tcell.Screen, x, y, w, h int, bg, border tcell.Color) {
	if w <= 0 || h <= 0 {
		return
	}

	bgStyle := tcell.StyleDefault.Background(bg)
	borderStyle := tcell.StyleDefault.Foreground(border).Background(bg)
	fillRect(screen, x, y, w, h, bgStyle)

	if w < 2 || h < 2 {
		return
	}

	screen.SetContent(x, y, '╭', nil, borderStyle)
	screen.SetContent(x+w-1, y, '╮', nil, borderStyle)
	screen.SetContent(x, y+h-1, '╰', nil, borderStyle)
	screen.SetContent(x+w-1, y+h-1, '╯', nil, borderStyle)

	for i := 1; i < w-1; i++ {
		screen.SetContent(x+i, y, '─', nil, borderStyle)
		screen.SetContent(x+i, y+h-1, '─', nil, borderStyle)
	}
	for j := 1; j < h-1; j++ {
		screen.SetContent(x, y+j, '│', nil, borderStyle)
		screen.SetContent(x+w-1, y+j, '│', nil, borderStyle)
	}
}

func fillRect(screen tcell.Screen, x, y, w, h int, style tcell.Style) {
	if w <= 0 || h <= 0 {
		return
	}

	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			screen.SetContent(x+col, y+row, ' ', nil, style)
		}
	}
}

func drawProgressBar(screen tcell.Screen, x, y, w int, progress float64, active, track tcell.Color) {
	if w <= 0 {
		return
	}

	trackStyle := tcell.StyleDefault.Foreground(track)
	activeStyle := tcell.StyleDefault.Foreground(active)
	for i := 0; i < w; i++ {
		screen.SetContent(x+i, y, '─', nil, trackStyle)
	}

	filled := int(math.Round(progress * float64(w-1)))
	if filled < 0 {
		filled = 0
	}
	if filled > w-1 {
		filled = w - 1
	}

	for i := 0; i <= filled; i++ {
		screen.SetContent(x+i, y, '━', nil, activeStyle)
	}
	screen.SetContent(x+filled, y, '●', nil, activeStyle)
}

func clickPositionToTime(x, barX, barW int, total time.Duration) time.Duration {
	if barW <= 1 || total <= 0 {
		return 0
	}

	ratio := float64(x-barX) / float64(barW-1)
	ratio = math.Max(0, math.Min(1, ratio))
	return time.Duration(ratio * float64(total))
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	wholeSeconds := int(d.Seconds() + 0.5)
	hours := wholeSeconds / 3600
	minutes := (wholeSeconds % 3600) / 60
	seconds := wholeSeconds % 60
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}

func formatSpeedLabel(speed float64) string {
	label := fmt.Sprintf("%.2f", speed)
	label = strings.TrimRight(label, "0")
	label = strings.TrimRight(label, ".")
	return label + "x"
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
