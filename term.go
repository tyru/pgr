package pgr

import (
	"fmt"
	"io"
)

func termSave(out io.Writer) {
	fmt.Fprint(out, "\x1b7")
}

func termRestore(out io.Writer) {
	fmt.Fprint(out, "\x1b8")
}

func termClearLine(out io.Writer) {
	fmt.Fprint(out, "\x1b\x5b2K")
}
