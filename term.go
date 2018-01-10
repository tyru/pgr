package pgr

import "fmt"

func termSave() {
	fmt.Print("\x1b7")
}

func termRestore() {
	fmt.Print("\x1b8")
}
