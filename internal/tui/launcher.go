package tui

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var ErrCancelled = errors.New("tui cancelled")

// ErrBackToMainMenu tells the caller to reopen the main launcher UI.
var ErrBackToMainMenu = errors.New("back to main menu")

type CommandSpec struct {
	Name             string
	Description      string
	Usage            string
	InputLabel       string
	InputPlaceholder string
}

var commandSpecs = []CommandSpec{
	{
		Name:             "play",
		Description:      "Play one audio file in the terminal.",
		Usage:            "echo play input.mp3",
		InputLabel:       "Audio file",
		InputPlaceholder: "example: song.mp3",
	},
	{
		Name:             "trim",
		Description:      "Trim a section of an audio file.",
		Usage:            "echo trim input.mp3 --start 00:30 --end 01:30 --out out.mp3",
		InputLabel:       "Arguments",
		InputPlaceholder: "input.mp3 --start 00:30 --end 01:30 --out out.mp3",
	},
	{
		Name:             "concat",
		Description:      "Concatenate multiple audio files into one.",
		Usage:            "echo concat file1.mp3 file2.mp3 --out out.mp3",
		InputLabel:       "Arguments",
		InputPlaceholder: "file1.mp3 file2.mp3 --out out.mp3",
	},
	{
		Name:             "extract",
		Description:      "Extract audio from a video file.",
		Usage:            "echo extract video.mp4 --out audio.mp3",
		InputLabel:       "Arguments",
		InputPlaceholder: "video.mp4 --out audio.mp3",
	},
	{
		Name:             "volume",
		Description:      "Increase or decrease audio volume.",
		Usage:            "echo volume input.mp3 --level 2.0 --out louder.mp3",
		InputLabel:       "Arguments",
		InputPlaceholder: "input.mp3 --level 2.0 --out louder.mp3",
	},
}

// CommandSpecs returns a copy of the launcher command metadata.
func CommandSpecs() []CommandSpec {
	specs := make([]CommandSpec, len(commandSpecs))
	copy(specs, commandSpecs)
	return specs
}

// Launch shows the interactive TUI launcher and returns argv-style arguments.
func Launch(binary string) ([]string, error) {
	app := tview.NewApplication().EnableMouse(true)
	pages := tview.NewPages()

	result := []string{binary}
	accepted := false

	title := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("[green::b]ECHO[-:-:-]\n[white]Terminal audio launcher[-]")
	title.SetBorder(true).SetTitle(" Welcome ")

	list := tview.NewList()
	list.ShowSecondaryText(true)
	list.SetBorder(true)
	list.SetTitle(" Commands ")
	list.SetHighlightFullLine(true)
	list.SetWrapAround(true)
	list.SetMainTextColor(tcell.ColorWhite)
	list.SetSecondaryTextColor(tcell.ColorGray)
	list.SetShortcutColor(tcell.ColorYellow)
	list.SetSelectedBackgroundColor(tcell.ColorGreen)
	list.SetSelectedTextColor(tcell.ColorBlack)

	detail := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true)
	detail.SetBorder(true).SetTitle(" Details ")

	for i, spec := range commandSpecs {
		list.AddItem(spec.Name, spec.Description, rune('1'+i), nil)
	}
	list.SetCurrentItem(0)
	detail.SetText(RenderDetails(commandSpecs[0]))

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetText("[gray]Use Up/Down, Enter, or click to choose. Esc exits.[-]")
	footer.SetBorder(true).SetTitle(" Help ")

	body := tview.NewFlex().SetDirection(tview.FlexColumn)
	body.AddItem(list, 0, 1, true)
	body.AddItem(detail, 0, 2, false)

	root := tview.NewFlex().SetDirection(tview.FlexRow)
	root.AddItem(title, 5, 0, false)
	root.AddItem(body, 0, 1, true)
	root.AddItem(footer, 3, 0, false)

	restoreMain := func() {
		if pages.HasPage("editor") {
			pages.RemovePage("editor")
		}
		if pages.HasPage("modal") {
			pages.RemovePage("modal")
		}
		pages.SwitchToPage("main")
		app.SetFocus(list)
	}

	showModal := func(titleText, message string, after func()) {
		if pages.HasPage("modal") {
			pages.RemovePage("modal")
		}

		modal := tview.NewModal().
			SetText(message).
			SetTextColor(tcell.ColorWhite).
			SetButtonTextColor(tcell.ColorBlack).
			SetButtonBackgroundColor(tcell.ColorGreen).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				pages.RemovePage("modal")
				if after != nil {
					after()
				}
			}).
			AddButtons([]string{"OK"})
		modal.SetBorder(true).SetTitle(" " + titleText + " ")

		pages.AddPage("modal", modal, true, true)
		app.SetFocus(modal)
	}

	runCommand := func(spec CommandSpec, input *tview.InputField) {
		text := strings.TrimSpace(input.GetText())

		if spec.Name == "play" {
			if text == "" {
				showModal("Error", "Please enter the audio file you want to play.", func() {
					app.SetFocus(input)
				})
				return
			}

			info, err := os.Stat(text)
			if err != nil {
				showModal("Error", fmt.Sprintf("Could not open the file: %v", err), func() {
					app.SetFocus(input)
				})
				return
			}

			if info.IsDir() {
				showModal("Error", "The path you entered is a directory, not a file.", func() {
					app.SetFocus(input)
				})
				return
			}

			result = []string{binary, spec.Name, text}
			accepted = true
			app.Stop()
			return
		}

		fields := strings.Fields(text)
		result = append([]string{binary, spec.Name}, fields...)
		accepted = true
		app.Stop()
	}

	openEditor := func(spec CommandSpec) {
		if pages.HasPage("editor") {
			pages.RemovePage("editor")
		}

		form := tview.NewForm()
		form.SetBorder(true)
		form.SetTitle(" " + strings.ToUpper(spec.Name) + " ")
		form.SetLabelColor(tcell.ColorLightCyan)
		form.SetFieldBackgroundColor(tcell.ColorBlack)
		form.SetFieldTextColor(tcell.ColorWhite)
		form.SetButtonBackgroundColor(tcell.ColorGreen)
		form.SetButtonTextColor(tcell.ColorBlack)
		form.SetButtonsAlign(tview.AlignRight)

		input := tview.NewInputField().
			SetLabel(spec.InputLabel + ":").
			SetPlaceholder(spec.InputPlaceholder).
			SetFieldWidth(50)

		input.SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				runCommand(spec, input)
			case tcell.KeyEscape:
				restoreMain()
			}
		})

		form.AddFormItem(input)
		form.AddButton("Run", func() {
			runCommand(spec, input)
		})
		form.AddButton("Back", func() {
			restoreMain()
		})
		form.SetCancelFunc(restoreMain)

		card := tview.NewFlex().SetDirection(tview.FlexColumn)
		card.AddItem(nil, 0, 1, false)
		card.AddItem(form, 70, 1, true)
		card.AddItem(nil, 0, 1, false)

		editor := tview.NewFlex().SetDirection(tview.FlexRow)
		editor.AddItem(nil, 0, 1, false)
		editor.AddItem(card, 14, 1, true)
		editor.AddItem(nil, 0, 1, false)

		pages.AddPage("editor", editor, true, true)
		app.SetFocus(input)
	}

	list.SetChangedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		if index >= 0 && index < len(commandSpecs) {
			detail.SetText(RenderDetails(commandSpecs[index]))
		}
	})
	list.SetSelectedFunc(func(index int, mainText, secondaryText string, shortcut rune) {
		if index >= 0 && index < len(commandSpecs) {
			openEditor(commandSpecs[index])
		}
	})
	list.SetDoneFunc(func() {
		app.Stop()
	})

	pages.AddPage("main", root, true, true)
	app.SetRoot(pages, true)
	app.SetFocus(list)

	if err := app.Run(); err != nil {
		return nil, err
	}

	if !accepted {
		return nil, ErrCancelled
	}

	return result, nil
}

// RenderDetails formats the selected command metadata for the details pane.
func RenderDetails(spec CommandSpec) string {
	return fmt.Sprintf(
		"[green::b]%s[-:-:-]\n\n[white]%s\n\n[gray]Usage:[white] %s\n[gray]Input:[white] %s\n",
		strings.ToUpper(spec.Name),
		spec.Description,
		spec.Usage,
		spec.InputPlaceholder,
	)
}
