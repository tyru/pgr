package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/tyru/pgr"
)

// This is the simplest example
func main() {
	poller := pgr.NewPoller(500 * time.Millisecond).
		Add(pgr.NewBarFunc(1, rocketAndMoai())).
		Add(pgr.NewBar(1, parseTemplate(`Press ENTER to exit...`)))
	ctx, cancel := context.WithCancel(context.Background())
	go poller.Show(ctx)

	os.Stdin.Read(make([]byte, 1))
	cancel()
	// Clear previous output + the entered line
	out := colorable.NewColorableStdout()
	fmt.Fprint(out, "\x1b\x5b3F"+"\x1b\x5b0K")
	// Print one more
	poller.Poll()
}

func rocketAndMoai() pgr.FormatFunc {
	i := 0
	return func(*pgr.Bar) string {
		i++
		return rocket(i) + moai(i)
	}
}

func rocket(i int) string {
	back := fmt.Sprintf("\x1b\x5b1D")
	i = i % 8
	if i < 4 {
		return strings.Repeat("  ", i) + "🚀" + strings.Repeat("  ", 4-i) + back
	} else if i%2 == 0 {
		return strings.Repeat("  ", 4) + "💥" + back
	} else {
		return strings.Repeat("  ", 4) + "🗯" + back
	}
}

func moai(i int) string {
	switch i % 5 {
	case 0:
		return "    🗿"
	case 1:
		return "   .🗿"
	case 2:
		return "  ｡.🗿"
	case 3:
		return " o｡.🗿"
	default:
		return "Oo｡.🗿"
	}
}

func parseTemplate(fmt string) *template.Template {
	return template.Must(template.New("moai").Parse(fmt))
}
