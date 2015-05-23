// DOS prompt only supports color attributes - no underline, no bold, no 256 colors
// TODO: what if user is using ConEmu?

package main

import "fmt"

func hello() {
	fmt.Println("this is windows")
}
