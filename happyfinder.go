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

// var termkey *TermboxEventWrapper

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

const pauseAfterKeypress = (500 * time.Millisecond)

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
	// fileChan := make(chan string, 1000)

	// go func() {
	// 	count := 0
	// 	// prefix := "brasilbrasilbrasilbrasilbrasil"
	// 	prefix := "brasilalemonalemonalemonalemonalemon"

	// 	for i := 0; i < 10000; i++ {
	// 		fileChan <- fmt.Sprintf("%s%d", prefix, count)
	// 		count += 1
	// 	}
	// }()

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

	// termkey = NewTermboxEventWrapper()

	modeline.Draw(&rview)
	cmdline.Draw(0, h-2, w)
	rview.SetSize(0, 0, w, h-2)
	// cmdline.Update(rview.GetSelected())
	// rview.CopyAll()
	// rview.Update()
	// rview.Draw()
	termbox.Flush()

	termboxEventChan := make(chan termbox.Event)
	// forceSortCh := make(chan bool, 10000)
	forceDrawCh := make(chan bool, 100)
	forceSortCh := make(chan bool, 100)
	// forceSortCh := make(chan bool)
	// resultCh := make(chan ResultSet, 1000)

	// go resultset.FilterManager(inputCh, resultCh)

	// throttle := time.NewTicker(time.Millisecond * 500)
	go func() {
		for {
			termboxEventChan <- termbox.PollEvent()
		}
	}()

	// dirty := false

	// go func() {
	// 	for {
	// 		filtered := <-resultCh
	// 		// dirty = false
	// 		rview.Update(filtered.results)
	// 		cmdline.Update(rview.GetSelected())

	// 		modeline.Draw(&rview)
	// 		cmdline.Draw(0, h-2, w)
	// 		rview.Draw()
	// 		termbox.Flush()
	// 	}
	// }()
	// var mutex = &sync.Mutex{}

	go func() {
		for {
			<-forceSortCh
			resultset.FlushQueue()
			filtered := resultset.Filter(global_lastkeypress, modeline.Contents())
			rview.Update(filtered.results)
			cmdline.Update(rview.GetSelected())
			forceDrawCh <- true
		}
	}()

	// func (rs *ResultSet) AsyncFilter(dirty <-chan bool, resultCh chan<- ResultSet) {
	// 	for {
	// 		<-dirty
	// 		result
	// 		res, _ := rs.Filter(when, userinput)
	// 		resultCh <- res
	// 	}
	// }

	idleTimer := time.NewTimer(1 * time.Hour)

	// metricsFilter = metrics.NewTimer()
	// metrics.Register("Filter", metricsFilter)

	// go metrics.Log(metrics.DefaultRegistry, 60e8, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	// Command name is:
	// os.Args[0]

	// var r string
	timeLastUser = time.Now().Add(-1 * time.Hour)
	timeLastFilter := time.Now()
	// dirty := false
	// ticker := time.NewTicker(time.Millisecond * 1000)

	for {
		select {
		case <-forceDrawCh:
			/* redraw */
		case <-idleTimer.C:
			// resultset.FlushQueue()
			// filtered := resultset.Filter(global_lastkeypress, modeline.Contents())
			// rview.Update(filtered.results)
			idleTimer = time.NewTimer(1 * time.Hour)
		// case <-ticker.C:
		// redraw
		// if dirty {
		// 	// go resultset.AsyncFilter(global_lastkeypress, modeline.Contents(), resultCh)
		// 	filtered, cancelled := resultset.Filter(global_lastkeypress, modeline.Contents())
		// 	if !cancelled {
		// 		rview.Update(filtered.results)
		// 		cmdline.Update(rview.GetSelected())
		// 		dirty = false
		// 	}
		// }
		case filename, ok := <-fileChan:
			if ok {
				// limit resorts

				// fmt.Println(filename, time.Since(timeLastUser))

				if time.Since(timeLastUser) > pauseAfterKeypress {

					modeline.Unpause()
					// resultset.Insert(filename)

					// dirty = true
					// resultset.AsyncFilter(global_lastkeypress, modeline.Contents(), resultCh)

					// filtered := resultset.Filter(global_lastkeypress, modeline.Contents())
					// rview.Update(filtered.results)
					// cmdline.Update(rview.GetSelected())
				} else {
					modeline.Pause()
					// resultset.Queue(filename)
				}

				resultset.Queue(filename)
				if !modeline.paused && time.Since(timeLastFilter) > (15*time.Millisecond) {
					forceSortCh <- true
					timeLastFilter = time.Now()
				}

			} else {
				// Last file received...
				resultset.FlushQueue()
				fileChan = nil
				forceSortCh <- true
			}

		case ev := <-termboxEventChan:
			if ev.Type == termbox.EventKey {
				timeLastUser = time.Now()
				global_lastkeypress = 0 //timeLastUser.UnixNano()
			}

			if fileChan != nil {
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
