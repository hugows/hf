package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	pauseAfterKeypress = (1500 * time.Millisecond)
	redrawPause        = 50 * time.Millisecond
)

var (
	global_lastkeypress int64
)

func main() {
	// defer profile.Start(profile.CPUProfile).Stop()

	opts, err := ParseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fi, err := os.Stat(opts.rootDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !fi.IsDir() {
		fmt.Println(opts.rootDir, "is NOT a folder")
		os.Exit(1)
	}

	stats := NewStats()

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc)

	fileset := new(ResultSet)

	var rview ResultsView
	windowWidth, windowHeight := termbox.Size()

	cmdline := NewCommandLine(opts.runCmd)
	modeline := NewModeline(opts.folderDisplay)

	idleTimer := time.NewTimer(1 * time.Hour)

	var fileCh <-chan string
	if opts.fakefiles > 0 {
		fileCh = walkFilesFake(2500)
	} else {
		fileCh = walkFiles(opts.rootDir)
	}

	termboxEventCh := make(chan termbox.Event)

	forceDrawCh := make(chan bool, 100)
	forceSortCh := make(chan bool, 100)

	timeLastUser := time.Now().Add(-1 * time.Hour)
	timeLastSort := time.Now()

	go func() {
		for {
			ev := termbox.PollEvent()
			stats.Inc("termboxEvent")
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
			stats.Inc("forceSort")
			filtered := fileset.Filter(global_lastkeypress, modeline.Contents())
			rview.Update(filtered.results)
			cmdline.Update(rview.GetMarkedOrSelected())
			forceDrawCh <- true
		}
	}()

	activeEditbox := modeline.input
	cmdline.SetActive(false)
	modeline.Draw(0, windowHeight-1, windowWidth, &rview, true)
	cmdline.Draw(0, windowHeight-2, windowWidth)
	rview.SetSize(0, 0, windowWidth, windowHeight-2)
	termbox.Flush()

	for {
		skipDraw := false

		select {
		case <-forceDrawCh:
			stats.Inc("forceDraw")
			rview.SelectFirst()

		case <-idleTimer.C:
			stats.Inc("idleTimer")
			idleTimer = time.NewTimer(1 * time.Hour)
			if !modeline.paused && fileCh == nil {
				forceSortCh <- true
			} else {
				idleTimer.Reset(redrawPause)
			}

		case filename, ok := <-fileCh:
			stats.Inc("fileCh")
			modeline.FlagPause(time.Since(timeLastUser) < pauseAfterKeypress)
			if ok {
				fileset.Insert(filename)
			} else {
				modeline.FlagLastFile()
				fileCh = nil
			}

			skipDraw = true
			if !modeline.paused && time.Since(timeLastSort) > redrawPause {
				forceSortCh <- true
				timeLastSort = time.Now()
			} else if !ok {
				idleTimer.Reset(redrawPause)
			}

		case ev := <-termboxEventCh:
			if fileCh == nil {
				modeline.FlagPause(false)
			} else {
				idleTimer.Reset(pauseAfterKeypress)
			}

			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc, termbox.KeyCtrlC:
					termbox.Close()
					if opts.debug {
						stats.Print()
					}
					return
				case termbox.KeyEnter:
					termbox.Close()
					runCmdWithArgs(opts.rootDir, cmdline.input.Contents(), false, cmdline.cmdargs)
					return
				case termbox.KeyCtrlT:
					err := rview.ToggleMarkAll()
					if err != nil {
						cmdline.ShowError(forceDrawCh, err)
					}
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

					cmdline.SetActive(activeEditbox == cmdline.input)
					// reset user edit to keep things simple...
					cmdline.Update(rview.GetMarkedOrSelected())

				case termbox.KeySpace:
					if activeEditbox == modeline.input {
						rview.ToggleMark()
						cmdline.Update(rview.GetMarkedOrSelected())
					} else {
						activeEditbox.InsertRune(' ') //? why ev.Ch failing??
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
			case termbox.EventResize:
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				windowHeight = ev.Height
				windowWidth = ev.Width
				rview.SetSize(0, 0, windowWidth, windowHeight-2)

			case termbox.EventError:
				termbox.Close()
				panic(ev.Err)
			}
		}

		if !skipDraw {
			modeline.Draw(0, windowHeight-1, windowWidth, &rview, activeEditbox == modeline.input)
			cmdline.Draw(0, windowHeight-2, windowWidth)
			rview.Draw()
			termbox.Flush()
			stats.Inc("flush")
		}
	}

}
