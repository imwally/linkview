package main

import (
	"fmt"
	"os/exec"
	"runtime"

	termbox "github.com/nsf/termbox-go"
)

type Terminal struct {
	Links         []Link
	Selected      int
	Width, Height int
}

const help string = "j/C-n: move down   k/C-p: move up   tab: show URL   return/C-o: open url   q: quit"

var (
	EventChan    = make(chan termbox.Event)
	KeyTab       = termbox.KeyTab
	KeyEnter     = termbox.KeyEnter
	KeyArrowUp   = termbox.KeyArrowUp
	KeyArrowDown = termbox.KeyArrowDown
	KeyCtrlP     = termbox.KeyCtrlP
	KeyCtrlN     = termbox.KeyCtrlN
	KeyCtrlO     = termbox.KeyCtrlO
	KeyCtrlC     = termbox.KeyCtrlC
)

func PollEvent() {
	EventChan <- termbox.PollEvent()
}

func NewTerminal(links *[]Link) *Terminal {
	term := Terminal{
		Links:    *links,
		Selected: 0,
		Width:    0,
		Height:   0,
	}

	return &term
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

func (t *Terminal) HandleEvent(e termbox.Event) (bool, error) {
	if e.Type == termbox.EventResize {
		t.SetSize()
		t.Render()
	}

	var err error
	if e.Type == termbox.EventKey {
		if e.Ch == 0 {
			switch e.Key {
			case KeyArrowDown:
				t.MoveSelection("down")
				t.Render()
			case KeyCtrlN:
				t.MoveSelection("down")
				t.Render()
			case KeyArrowUp:
				t.MoveSelection("up")
				t.Render()
			case KeyCtrlP:
				t.MoveSelection("up")
				t.Render()
			case KeyTab:
				t.ShowFullLink()
			case KeyEnter:
				err = t.Select()
			case KeyCtrlO:
				err = t.Select()
			case KeyCtrlC:
				return true, nil
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
				return true, nil
			}
		}
	}

	return false, err
}

func (t *Terminal) Println(x int, y int, s string) {
	for col, char := range s {
		termbox.SetCell(col+x, y, char, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (t *Terminal) ShowFullLink() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	t.Println(0, 0, help)

	url := t.Links[t.Selected].URL
	row := 2
	col := 0
	for _, char := range url {
		if col >= t.Width {
			row++
			col = 0
		}
		termbox.SetCell(col, row, char, termbox.ColorDefault, termbox.ColorDefault)
		col++
	}

	termbox.Flush()
}

func (t *Terminal) Render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	url := t.Links[t.Selected].URL
	t.Println(0, 0, help)
	t.Println(0, 2, url)

	var start int
	offset := t.Selected - t.Height + 6
	if t.Selected > t.Height-6 {
		start = offset
	}

	links := t.Links
	for i := start; i < len(links); i++ {
		if t.Selected > t.Height-6 {
			t.Println(0, t.Height-2, "->")
			t.Println(3, i+4-offset, links[i].Text)
		} else {
			t.Println(0, t.Selected+4, "->")
			t.Println(3, i+4, links[i].Text)
		}
	}

	termbox.Flush()
}

func (t *Terminal) SetSize() {
	t.Width, t.Height = termbox.Size()
}

func (t *Terminal) MoveSelection(direction string) {
	switch direction {
	case "up":
		t.Selected--
	case "down":
		t.Selected++
	}

	if t.Selected >= len(t.Links) {
		t.Selected = len(t.Links) - 1
	}

	if t.Selected < 0 {
		t.Selected = 0
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
		err = fmt.Errorf("can't open browser: unsupported platform")
	}

	if err != nil {
		return err
	}

	return nil
}
