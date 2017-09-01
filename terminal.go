package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"

	termbox "github.com/nsf/termbox-go"
)

type Terminal struct {
	Links    []Link
	Selected int
}

var (
	EventChan    = make(chan termbox.Event)
	KeyTab       = termbox.KeyTab
	KeyEnter     = termbox.KeyEnter
	KeyArrowUp   = termbox.KeyArrowUp
	KeyArrowDown = termbox.KeyArrowDown
)

func NewTerminal(links *[]Link) *Terminal {
	term := Terminal{
		*links,
		0,
	}

	return &term
}

func PollEvent() {
	event := termbox.PollEvent()
	if event.Type == termbox.EventKey {
		EventChan <- event
	}
}

func (t *Terminal) Start() error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	return nil
}

func (t *Terminal) Close() {
	termbox.Close()
}

func (t *Terminal) HandleEvent(e termbox.Event) bool {
	if e.Ch == 0 {
		switch e.Key {
		case KeyArrowDown:
			t.MoveSelection("down")
			t.Render()
		case KeyArrowUp:
			t.MoveSelection("up")
			t.Render()
		case KeyEnter:
			t.Select()
		}
	} else {
		switch e.Ch {
		case 'j':
			t.MoveSelection("down")
			t.Render()
		case 'k':
			t.MoveSelection("up")
			t.Render()
		case 'q':
			return true
		}
	}

	return false
}

func (t *Terminal) Render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	help := "j: move down   k: move up   return: open url   q: quit"
	for i, character := range help {
		termbox.SetCell(i, 0, character, termbox.ColorDefault, termbox.ColorDefault)
	}

	for col, character := range t.Links[t.Selected].URL {
		termbox.SetCell(col+2, 2, character, termbox.ColorDefault, termbox.ColorDefault)
	}

	for row, link := range t.Links {
		if row == t.Selected {
			termbox.SetCell(0, row+4, '→', termbox.ColorDefault, termbox.ColorDefault)
			termbox.SetCell(1, row+4, ' ', termbox.ColorDefault, termbox.ColorDefault)
		}
		for col, character := range link.Text {
			termbox.SetCell(col+2, row+4, character, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	termbox.Flush()
}

func (t *Terminal) MoveSelection(direction string) {
	switch direction {
	case "up":
		t.Selected--
	case "down":
		t.Selected++
	}

	if t.Selected >= len(t.Links) {
		t.Selected = 0
	}
	if t.Selected < 0 {
		t.Selected = len(t.Links) - 1
	}
}

func (t *Terminal) Select() error {
	var err error
	url := t.Links[t.Selected].URL

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
