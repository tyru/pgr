// +build ignore

package main

import (
	"context"
	"fmt"
	"math/rand"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/tyru/pgr"
)

func main() {
	b1 := makeBar("b1", "fgCyan")
	b2 := makeBar("b2", "fgYellow")
	b3 := makeBar("b3", "fgRed")
	b4 := makeBar("b4", "fgBlue")
	b5 := makeBar("b5", "fgGreen")

	poller := pgr.NewPoller(300 * time.Millisecond).
		Add(b1).
		Add(b2).
		Add(b3).
		Add(b4).
		Add(b5).
		SetSorter(func(bars []*pgr.Bar, i, j int) bool {
			return bars[i].Current() >= bars[j].Current()
		})

	now := time.Now().Unix()
	go randIncBy(b1, 300*time.Millisecond, now+1)
	go randIncBy(b2, 300*time.Millisecond, now+2)
	go randIncBy(b3, 300*time.Millisecond, now+3)
	go randIncBy(b4, 300*time.Millisecond, now+4)
	go randIncBy(b5, 300*time.Millisecond, now+5)

	poller.Show(context.Background())
}

func makeBar(name, color string) *pgr.Bar {
	prefix := fmt.Sprintf("({{ %s %q }})", color, name)
	funcMap := funcMap()
	return pgr.NewBar(100, parseTemplate(funcMap, prefix+` {{ bar . "[" "=" ">" " " "]" 70 }}`)).
		OnFinish(parseTemplate(funcMap, prefix+` Finished!`))
}

func parseTemplate(funcMap template.FuncMap, format string) *template.Template {
	return template.Must(pgr.NewTemplate().Funcs(funcMap).Parse(format))
}

func funcMap() template.FuncMap {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	return template.FuncMap{
		"fgCyan": func(s string) string {
			return cyan(s)
		},
		"fgYellow": func(s string) string {
			return yellow(s)
		},
		"fgRed": func(s string) string {
			return red(s)
		},
		"fgBlue": func(s string) string {
			return blue(s)
		},
		"fgGreen": func(s string) string {
			return green(s)
		},
	}
}

func randIncBy(p *pgr.Bar, d time.Duration, seed int64) {
	r := rand.New(rand.NewSource(seed))
	for {
		p.Add(int64(r.Int() % 5))
		time.Sleep(d)
	}
}
