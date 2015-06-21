package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/nsf/termbox-go"
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

// strings.Replace(tw.Text, " ", "+", -1)

func runCmdWithArgs(f string) {
	cmd := exec.Command(*cmd, f)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

const pauseAfterKeypress = (1000 * time.Millisecond)

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

	fileChan := walkFiles(getRoot())

	resultset := new(ResultSet)

	// for filename := range fileChan {
	// results.Insert(<-fileChan)
	// }

	// 	a, b := score(filexname, flag.Arg(1))
	// 	if a >= 0 && a < 100 {

	// 		if first {
	// 			runcmd := exec.Command(*cmd, filename)
	// 			runcmd.Stdin = os.Stdin
	// 			runcmd.Stdout = os.Stdout
	// 			err := runcmd.Run()
	// 			if err != nil {
	// 				log.Fatal(err)
	// 			}
	// 			first = false
	// 		}

	// 		fmt.Printf("%30s %4d %v\n", filename, a, b)
	// 	}
	// }

	var timeLastUser time.Time
	// resultsQueue := make([]string, 0, 100)
	w, h := termbox.Size()
	modeline := NewModeline(0, h-1, w)
	cmdline := new(CommandLine)

	modeline.Draw(&rview)
	cmdline.Draw(0, h-2, w)
	rview.SetSize(0, 0, w, h-2)
	// cmdline.Update(rview.GetSelected())
	// rview.CopyAll()
	// rview.Update()
	// rview.Draw()
	termbox.Flush()

	termboxEventChan := make(chan termbox.Event)
	resultCh := make(chan ResultSet)
	quit := make(chan bool, 1)

	go func() {
		for {
			termboxEventChan <- termbox.PollEvent()
		}
	}()

	timer := time.NewTimer(1 * time.Hour)

	// Command name is:
	// os.Args[0]

	// var r string
	timeLastUser = time.Now().Add(-1 * time.Hour)

	for {
		select {
		case filtered := <-resultCh:
			rview.Update(filtered.results)
			cmdline.Update(rview.GetSelected())
		case <-timer.C:
			resultset.FlushQueue()
			resultset.AsyncFilter(modeline.Contents(), resultCh, quit)
			timer = time.NewTimer(1 * time.Hour)
		case filename, ok := <-fileChan:
			if ok {
				if time.Since(timeLastUser) > pauseAfterKeypress {
					modeline.Unpause()
					resultset.Insert(filename)
					resultset.AsyncFilter(modeline.Contents(), resultCh, quit)
				} else {
					modeline.Pause()
					resultset.Queue(filename)
				}
			} else {
				fileChan = nil
			}

		case ev := <-termboxEventChan:
			if ev.Type == termbox.EventKey {
				timeLastUser = time.Now()
			}

			if fileChan != nil {
				timer.Reset(pauseAfterKeypress)
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
					resultset.AsyncFilter(modeline.Contents(), resultCh, quit)
				case termbox.KeyDelete, termbox.KeyCtrlD:
					modeline.input.DeleteRuneForward()
					resultset.AsyncFilter(modeline.Contents(), resultCh, quit)
				case termbox.KeySpace:
					rview.ToggleMark()
				case termbox.KeyCtrlK:
					modeline.input.DeleteTheRestOfTheLine()
					resultset.AsyncFilter(modeline.Contents(), resultCh, quit)
				case termbox.KeyHome, termbox.KeyCtrlA:
					modeline.input.MoveCursorToBeginningOfTheLine()
				case termbox.KeyEnd, termbox.KeyCtrlE:
					modeline.input.MoveCursorToEndOfTheLine()
				default:
					if ev.Ch != 0 {
						modeline.input.InsertRune(ev.Ch)
						resultset.AsyncFilter(modeline.Contents(), resultCh, quit)
					}
				}
			case termbox.EventError:
				panic(ev.Err)
			}
		}

		modeline.Draw(&rview)
		cmdline.Draw(0, h-2, w)
		rview.Draw()
		termbox.Flush()
	}

}
