package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

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

// const modeline_width = 30

func redraw_all(modeline Modeline) {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	w, h := termbox.Size()
	// midy := h / 2
	// midx := (w - modeline_width) / 2

	// unicode box drawing chars around the edit box
	// termbox.SetCell(midx-1, midy, '│', coldef, coldef)
	// termbox.SetCell(midx+modeline_width, midy, '│', coldef, coldef)
	// termbox.SetCell(midx-1, midy-1, '┌', coldef, coldef)
	// termbox.SetCell(midx-1, midy+1, '└', coldef, coldef)
	// termbox.SetCell(midx+modeline_width, midy-1, '┐', coldef, coldef)
	// termbox.SetCell(midx+modeline_width, midy+1, '┘', coldef, coldef)
	// fill(midx, midy-1, modeline_width, 1, termbox.Cell{Ch: '─'})

	// fill(0, h-2, modeline_width+2, 1, termbox.Cell{Ch: '─'})

	termbox.SetCell(0, h-1, '>', coldef, coldef)
	modeline.Draw(2, h-1, w-2, 1)
	termbox.SetCursor(2+modeline.CursorX(), h-1)

	//tbprint(0, h-1, coldef, coldef, "Press ESC to quit")
}

func runCmdWithArgs(f string) {
	// fmt.Println(*cmd, f)
	cmd := exec.Command(*cmd, f)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	var modeline Modeline
	var statusline Statusline
	var results Results

	_, err := os.Open(getRoot())
	if err != nil {
		panic(err)
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc)

	files := walkFiles(getRoot())

	for filename := range files {
		results.Insert(filename)
	}

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

	w, h := termbox.Size()
	redraw_all(modeline)
	statusline.Draw(0, h-2, w, &results)
	results.SetSize(0, 0, w, h-2)
	results.CopyAll()
	results.Draw()
	termbox.Flush()

	// Command name is:
	// os.Args[0]

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				termbox.Close()
				return
			case termbox.KeyEnter:
				termbox.Close()
				// runCmdWithArgs(results.FormatSelected())
				return
			case termbox.KeyCtrlT:
				results.ToggleMarkAll()
			case termbox.KeyArrowUp, termbox.KeyCtrlP:
				results.SelectPrevious()
			case termbox.KeyArrowDown, termbox.KeyCtrlN:
				results.SelectNext()
			case termbox.KeyArrowLeft, termbox.KeyCtrlB:
				modeline.MoveCursorOneRuneBackward()
			case termbox.KeyArrowRight, termbox.KeyCtrlF:
				modeline.MoveCursorOneRuneForward()
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				modeline.DeleteRuneBackward()
				results.Filter(modeline.Contents())
			case termbox.KeyDelete, termbox.KeyCtrlD:
				modeline.DeleteRuneForward()
				results.Filter(modeline.Contents())
			case termbox.KeySpace:
				results.ToggleMark()
			case termbox.KeyCtrlK:
				modeline.DeleteTheRestOfTheLine()
				results.Filter(modeline.Contents())
			case termbox.KeyHome, termbox.KeyCtrlA:
				modeline.MoveCursorToBeginningOfTheLine()
			case termbox.KeyEnd, termbox.KeyCtrlE:
				modeline.MoveCursorToEndOfTheLine()
			default:
				if ev.Ch != 0 {
					modeline.InsertRune(ev.Ch)
					results.Filter(modeline.Contents())
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redraw_all(modeline)
		statusline.Draw(0, h-2, w, &results)

		results.Draw()
		termbox.Flush()
	}

}
