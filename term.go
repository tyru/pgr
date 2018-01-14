package pgr

import (
	"fmt"
	"io"
)

func termPrevLine(out io.Writer, n int) {
	fmt.Fprintf(out, "\x1b\x5b%dF", n)
}

func termClearLine(out io.Writer) {
	fmt.Fprint(out, "\x1b\x5b2K")
}
