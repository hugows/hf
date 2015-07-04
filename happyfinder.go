package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	pauseAfterKeypress = (1500 * time.Millisecond)
	redrawPause        = 30 * time.Millisecond
)

var (
	global_lastkeypress int64
)

func getRoot() string {
	if len(flag.Args()) == 0 {
		return "."
	} else {
		return flag.Arg(0)
	}
}

// hf --cmd=emacs ~/go/src/github.com/hugows/ happy
var cmd = flag.String("cmd", "vim", "command to run")

func main() {
	flag.Parse()

	var rview ResultsView

	root := getRoot()
	fi, err := os.Stat(root)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !fi.IsDir() {
		fmt.Println(root, "is NOT a folder")
		return
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc)

	fileset := new(ResultSet)

	w, h := termbox.Size()
	modeline := NewModeline(0, h-1, w)
	cmdline := new(CommandLine)

	idleTimer := time.NewTimer(1 * time.Hour)

	fileCh := walkFiles(getRoot())
	termboxEventCh := make(chan termbox.Event)

	forceDrawCh := make(chan bool, 100)
	forceSortCh := make(chan bool, 100)

	timeLastUser := time.Now().Add(-1 * time.Hour)
	timeLastFilter := time.Now()

	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey {
				timeLastUser = time.Now()
				global_lastkeypress = timeLastUser.UnixNano()
			}
			termboxEventCh <- ev
		}
	}()

	go func() {
		for {
			<-forceSortCh
			filtered := fileset.Filter(global_lastkeypress, modeline.Contents())
			rview.Update(filtered.results)
			cmdline.Update(rview.GetSelected())
			forceDrawCh <- true
		}
	}()

	// Command name is:
	// os.Args[0]

	modeline.Draw(&rview)
	cmdline.Draw(0, h-2, w)
	rview.SetSize(0, 0, w, h-2)
	termbox.Flush()

	for {
		select {
		case <-forceDrawCh:
			/* redraw */

		case <-idleTimer.C:
			idleTimer = time.NewTimer(1 * time.Hour)
			if !modeline.paused {
				modeline.LastFile()
				fileCh = nil
				forceSortCh <- true
			}

		case filename, ok := <-fileCh:
			if time.Since(timeLastUser) > pauseAfterKeypress {
				modeline.Unpause()
			} else {
				modeline.Pause()
			}

			if ok {
				fileset.Insert(filename)
			}

			if !modeline.paused && time.Since(timeLastFilter) > redrawPause {
				forceSortCh <- true
				timeLastFilter = time.Now()

				if !ok {
					modeline.LastFile()
					fileCh = nil
				}
			} else if !ok {
				idleTimer.Reset(redrawPause)
				fileCh = nil
			}

		case ev := <-termboxEventCh:
			if fileCh != nil {
				idleTimer.Reset(pauseAfterKeypress)
			} else {
				modeline.Unpause()
			}

			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc, termbox.KeyCtrlC:
					termbox.Close()
					return
				case termbox.KeyEnter:
					termbox.Close()
					// runCmdWithArgs(rview.FormatSelected())
					return
				case termbox.KeyCtrlT:
					rview.ToggleMarkAll()
				case termbox.KeyArrowUp, termbox.KeyCtrlP:
					cmdline.Update(rview.SelectPrevious())
				case termbox.KeyArrowDown, termbox.KeyCtrlN:
					cmdline.Update(rview.SelectNext())
				case termbox.KeyArrowLeft, termbox.KeyCtrlB:
					modeline.input.MoveCursorOneRuneBackward()
				case termbox.KeyArrowRight, termbox.KeyCtrlF:
					modeline.input.MoveCursorOneRuneForward()
				case termbox.KeyBackspace, termbox.KeyBackspace2:
					modeline.input.DeleteRuneBackward()
					forceSortCh <- true
				case termbox.KeyDelete, termbox.KeyCtrlD:
					modeline.input.DeleteRuneForward()
					forceSortCh <- true
				case termbox.KeySpace:
					rview.ToggleMark()
				case termbox.KeyCtrlK:
					modeline.input.DeleteTheRestOfTheLine()
					forceSortCh <- true
				case termbox.KeyHome, termbox.KeyCtrlA:
					modeline.input.MoveCursorToBeginningOfTheLine()
				case termbox.KeyEnd, termbox.KeyCtrlE:
					modeline.input.MoveCursorToEndOfTheLine()
				default:
					if ev.Ch != 0 {
						modeline.input.InsertRune(ev.Ch)
						forceSortCh <- true
					}
				}
			case termbox.EventError:
				panic(ev.Err)
			}

			// fmt.Println(modeline.Contents())
		}

		modeline.Draw(&rview)
		cmdline.Draw(0, h-2, w)
		rview.Draw()
		termbox.Flush()
	}

}
