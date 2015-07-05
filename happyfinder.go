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
	cmdline := NewCommandLine(0, h-2, w, "vim")
	modeline := NewModeline(0, h-1, w)

	idleTimer := time.NewTimer(1 * time.Hour)

	fileCh := walkFiles(getRoot())
	// fileCh := walkFilesFake(2500)
	termboxEventCh := make(chan termbox.Event)

	forceDrawCh := make(chan bool, 100)
	forceSortCh := make(chan bool, 100)

	timeLastUser := time.Now().Add(-1 * time.Hour)
	timeLastSort := time.Now()

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
			cmdline.Update(rview.GetMarkedOrSelected())
			forceDrawCh <- true
		}
	}()

	// Command name is:
	// os.Args[0]

	activeEditbox := modeline.input

	modeline.Draw(&rview, true)
	cmdline.Draw(0, h-2, w, false)
	rview.SetSize(0, 0, w, h-2)
	termbox.Flush()

	for {
		select {
		case <-forceDrawCh:
			// bug
			// rview.SelectFirst()
			/* redraw */

		case <-idleTimer.C:
			idleTimer = time.NewTimer(1 * time.Hour)
			if !modeline.paused && fileCh == nil {
				forceSortCh <- true
			} else {
				idleTimer.Reset(redrawPause)
			}

		case filename, ok := <-fileCh:
			modeline.FlagPause(time.Since(timeLastUser) < pauseAfterKeypress)

			if ok {
				fileset.Insert(filename)
			} else {
				modeline.FlagLastFile()
				fileCh = nil
			}

			if !modeline.paused && time.Since(timeLastSort) > redrawPause {
				forceSortCh <- true
				timeLastSort = time.Now()
			} else if !ok {
				idleTimer.Reset(redrawPause)
			}

		case ev := <-termboxEventCh:
			idleTimer.Reset(pauseAfterKeypress)
			if fileCh == nil {
				modeline.FlagPause(false)
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
					cmdline.Update(rview.GetMarkedOrSelected())
				case termbox.KeyArrowUp, termbox.KeyCtrlP:
					rview.SelectPrevious()
					cmdline.Update(rview.GetMarkedOrSelected())
				case termbox.KeyArrowDown, termbox.KeyCtrlN:
					rview.SelectNext()
					cmdline.Update(rview.GetMarkedOrSelected())
				case termbox.KeyArrowLeft, termbox.KeyCtrlB:
					activeEditbox.MoveCursorOneRuneBackward()
				case termbox.KeyArrowRight, termbox.KeyCtrlF:
					activeEditbox.MoveCursorOneRuneForward()
				case termbox.KeyBackspace, termbox.KeyBackspace2:
					activeEditbox.DeleteRuneBackward()
					if activeEditbox == modeline.input {
						forceSortCh <- true
					}
				case termbox.KeyDelete, termbox.KeyCtrlD:
					activeEditbox.DeleteRuneForward()
					if activeEditbox == modeline.input {
						forceSortCh <- true
					}
				case termbox.KeyTab:
					if activeEditbox == modeline.input {
						activeEditbox = cmdline.input
					} else {
						activeEditbox = modeline.input
					}
				case termbox.KeySpace:
					if activeEditbox == modeline.input {
						rview.ToggleMark()
						cmdline.Update(rview.GetMarkedOrSelected())
					} else {
						activeEditbox.InsertRune(ev.Ch)
					}
				case termbox.KeyCtrlK:
					activeEditbox.DeleteTheRestOfTheLine()
					if activeEditbox == modeline.input {
						forceSortCh <- true
					}
				case termbox.KeyHome, termbox.KeyCtrlA:
					activeEditbox.MoveCursorToBeginningOfTheLine()
				case termbox.KeyEnd, termbox.KeyCtrlE:
					activeEditbox.MoveCursorToEndOfTheLine()
				default:
					if ev.Ch != 0 {
						activeEditbox.InsertRune(ev.Ch)
						if activeEditbox == modeline.input {
							forceSortCh <- true
						}
					}
				}
			case termbox.EventError:
				panic(ev.Err)
			}
		}

		modeline.Draw(&rview, activeEditbox == modeline.input)
		cmdline.Draw(0, h-2, w, activeEditbox == cmdline.input)
		rview.Draw()
		termbox.Flush()
	}

}
