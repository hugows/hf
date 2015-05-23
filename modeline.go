package main

import (
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

type Modeline struct {
	text           []byte
	line_voffset   int
	cursor_boffset int // cursor offset in bytes
	cursor_voffset int // visual cursor offset in termbox cells
	cursor_coffset int // cursor offset in unicode code points
}

// Draws the Modeline in the given location, 'h' is not used at the moment
func (m *Modeline) Draw(x, y, w, h int) {
	m.AdjustVOffset(w)

	const coldef = termbox.ColorDefault
	fill(x, y, w, h, termbox.Cell{Ch: ' '})

	t := m.text
	lx := 0

	for {
		rx := lx - m.line_voffset
		if len(t) == 0 {
			break
		}

		if rx >= w {
			termbox.SetCell(x+w-1, y, '→',
				coldef, coldef)
			break
		}

		r, size := utf8.DecodeRune(t)
		if rx >= 0 {
			termbox.SetCell(x+rx, y, r, coldef, coldef)
		}
		lx += 1
		t = t[size:]
	}

	if m.line_voffset != 0 {
		termbox.SetCell(x, y, '←', coldef, coldef)
	}
}

// Adjusts line visual offset to a proper value depending on width
func (m *Modeline) AdjustVOffset(width int) {
	ht := 0

	threshold := width - 1
	if m.line_voffset != 0 {
		threshold = width - ht
	}
	if m.cursor_voffset-m.line_voffset >= threshold {
		m.line_voffset = m.cursor_voffset + (ht - width + 1)
	}

	if m.line_voffset != 0 && m.cursor_voffset-m.line_voffset < ht {
		m.line_voffset = m.cursor_voffset - ht
		if m.line_voffset < 0 {
			m.line_voffset = 0
		}
	}
}

func (m *Modeline) Contents() string {
	return string(m.text)
}

func (m *Modeline) MoveCursorTo(boffset int) {
	m.cursor_boffset = boffset
	m.cursor_voffset, m.cursor_coffset = voffset_coffset(m.text, boffset)
}

func (m *Modeline) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(m.text[m.cursor_boffset:])
}

func (m *Modeline) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(m.text[:m.cursor_boffset])
}

func (m *Modeline) MoveCursorOneRuneBackward() {
	if m.cursor_boffset == 0 {
		return
	}
	_, size := m.RuneBeforeCursor()
	m.MoveCursorTo(m.cursor_boffset - size)
}

func (m *Modeline) MoveCursorOneRuneForward() {
	if m.cursor_boffset == len(m.text) {
		return
	}
	_, size := m.RuneUnderCursor()
	m.MoveCursorTo(m.cursor_boffset + size)
}

func (m *Modeline) MoveCursorToBeginningOfTheLine() {
	m.MoveCursorTo(0)
}

func (m *Modeline) MoveCursorToEndOfTheLine() {
	m.MoveCursorTo(len(m.text))
}

func (m *Modeline) DeleteRuneBackward() {
	if m.cursor_boffset == 0 {
		return
	}

	m.MoveCursorOneRuneBackward()
	_, size := m.RuneUnderCursor()
	m.text = byte_slice_remove(m.text, m.cursor_boffset, m.cursor_boffset+size)
}

func (m *Modeline) DeleteRuneForward() {
	if m.cursor_boffset == len(m.text) {
		return
	}
	_, size := m.RuneUnderCursor()
	m.text = byte_slice_remove(m.text, m.cursor_boffset, m.cursor_boffset+size)
}

func (m *Modeline) DeleteTheRestOfTheLine() {
	m.text = m.text[:m.cursor_boffset]
}

func (m *Modeline) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	m.text = byte_slice_insert(m.text, m.cursor_boffset, buf[:n])
	m.MoveCursorOneRuneForward()
}

// Please, keep in mind that cursor depends on the value of line_voffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (m *Modeline) CursorX() int {
	return m.cursor_voffset - m.line_voffset
}

// func main() {
// 	err := termbox.Init()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer termbox.Close()
// 	termbox.SetInputMode(termbox.InputEsc)

// 	redraw_all()
// mainloop:
// 	for {
// 		switch ev := termbox.PollEvent(); ev.Type {
// 		case termbox.EventKey:
// 			switch ev.Key {
// 			case termbox.KeyEsc:
// 				break mainloop
// 			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
// 				edit_box.MoveCursorOneRuneBackward()
// 			case termbox.KeyArrowRight, termbox.KeyCtrlF:
// 				edit_box.MoveCursorOneRuneForward()
// 			case termbox.KeyBackspace, termbox.KeyBackspace2:
// 				edit_box.DeleteRuneBackward()
// 			case termbox.KeyDelete, termbox.KeyCtrlD:
// 				edit_box.DeleteRuneForward()
// 			case termbox.KeyTab:
// 				edit_box.InsertRune('\t')
// 			case termbox.KeySpace:
// 				edit_box.InsertRune(' ')
// 			case termbox.KeyCtrlK:
// 				edit_box.DeleteTheRestOfTheLine()
// 			case termbox.KeyHome, termbox.KeyCtrlA:
// 				edit_box.MoveCursorToBeginningOfTheLine()
// 			case termbox.KeyEnd, termbox.KeyCtrlE:
// 				edit_box.MoveCursorToEndOfTheLine()
// 			default:
// 				if ev.Ch != 0 {
// 					edit_box.InsertRune(ev.Ch)
// 				}
// 			}
// 		case termbox.EventError:
// 			panic(ev.Err)
// 		}
// 		redraw_all()
// 	}
// }
