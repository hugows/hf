package main

import (
	"unicode/utf8"

	"github.com/nsf/termbox-go"
)

type Editbox struct {
	text           []byte
	line_voffset   int
	cursor_boffset int // cursor offset in bytes
	cursor_voffset int // visual cursor offset in termbox cells
	cursor_coffset int // cursor offset in unicode code points

	// colors
	fg termbox.Attribute
	bg termbox.Attribute
}

// Draws the Editbox in the given location
func (e *Editbox) Draw(x, y, w int) {
	e.AdjustVOffset(w)

	tclear(x, y, w, 1)

	t := e.text
	lx := 0

	for {
		rx := lx - e.line_voffset
		if len(t) == 0 {
			break
		}

		if rx >= w {
			termbox.SetCell(x+w-1, y, '→', e.fg, e.bg)
			break
		}

		r, size := utf8.DecodeRune(t)
		if rx >= 0 {
			termbox.SetCell(x+rx, y, r, e.fg, e.bg)
		}
		lx += 1
		t = t[size:]
	}

	if e.line_voffset != 0 {
		termbox.SetCell(x, y, '←', e.fg, e.bg)
	}
}

// Adjusts line visual offset to a proper value depending on width
func (m *Editbox) AdjustVOffset(width int) {
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

func (m *Editbox) Contents() string {
	return string(m.text)
}

func (m *Editbox) MoveCursorTo(boffset int) {
	m.cursor_boffset = boffset
	m.cursor_voffset, m.cursor_coffset = voffset_coffset(m.text, boffset)
}

func (m *Editbox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(m.text[m.cursor_boffset:])
}

func (m *Editbox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(m.text[:m.cursor_boffset])
}

func (m *Editbox) MoveCursorOneRuneBackward() {
	if m.cursor_boffset == 0 {
		return
	}
	_, size := m.RuneBeforeCursor()
	m.MoveCursorTo(m.cursor_boffset - size)
}

func (m *Editbox) MoveCursorOneRuneForward() {
	if m.cursor_boffset == len(m.text) {
		return
	}
	_, size := m.RuneUnderCursor()
	m.MoveCursorTo(m.cursor_boffset + size)
}

func (m *Editbox) MoveCursorToBeginningOfTheLine() {
	m.MoveCursorTo(0)
}

func (m *Editbox) MoveCursorToEndOfTheLine() {
	m.MoveCursorTo(len(m.text))
}

func (m *Editbox) DeleteRuneBackward() {
	if m.cursor_boffset == 0 {
		return
	}

	m.MoveCursorOneRuneBackward()
	_, size := m.RuneUnderCursor()
	m.text = byte_slice_remove(m.text, m.cursor_boffset, m.cursor_boffset+size)
}

func (m *Editbox) DeleteRuneForward() {
	if m.cursor_boffset == len(m.text) {
		return
	}
	_, size := m.RuneUnderCursor()
	m.text = byte_slice_remove(m.text, m.cursor_boffset, m.cursor_boffset+size)
}

func (m *Editbox) DeleteTheRestOfTheLine() {
	m.text = m.text[:m.cursor_boffset]
}

func (m *Editbox) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	m.text = byte_slice_insert(m.text, m.cursor_boffset, buf[:n])
	m.MoveCursorOneRuneForward()
}

// Please, keep in mind that cursor depends on the value of line_voffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (m *Editbox) CursorX() int {
	return m.cursor_voffset - m.line_voffset
}
