// +build ignore

package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"github.com/tyru/pgr"
)

func main() {
	maxWidth := runewidth.StringWidth("ﾒｲｶｲｴｸﾘﾌﾟｽ")
	b1 := makeBar("ﾋﾞﾑﾃｲｵｰ", "fgCyan", maxWidth)
	b2 := makeBar("ｲｰﾏｸｽ", "fgYellow", maxWidth)
	b3 := makeBar("ﾒｲｶｲｴｸﾘﾌﾟｽ", "fgRed", maxWidth)
	b4 := makeBar("ｱﾄﾑｻﾞｺｰﾄﾞ", "fgBlue", maxWidth)
	b5 := makeBar("ｻﾌﾞﾗｲﾑﾃｷｽﾄ", "fgGreen", maxWidth)

	poller := pgr.NewPoller(250 * time.Millisecond).
		Add(b1).
		Add(b2).
		Add(b3).
		Add(b4).
		Add(b5).
		SetSorter(func(bars []*pgr.Bar, i, j int) bool {
			return bars[i].Current() < bars[j].Current()
		})

	go randIncBy(b1, 500*time.Millisecond, 1)
	go randIncBy(b2, 500*time.Millisecond, 2)
	go randIncBy(b3, 500*time.Millisecond, 3)
	go randIncBy(b4, 500*time.Millisecond, 4)
	go randIncBy(b5, 500*time.Millisecond, 5)

	poller.Show(context.Background())
}

func makeBar(name, color string, width int) *pgr.Bar {
	prefix := fmt.Sprintf("({{ %s %q }})", color, name)
	pad := strings.Repeat(" ", width-runewidth.StringWidth(name))
	return pgr.NewBar(100, parseTemplate(prefix+pad+` {{ bar . "[" "=" ">" " " "]" 70 }}`)).
		OnFinish(parseTemplate(prefix + pad + ` Goal!`))
}

func parseTemplate(format string) *template.Template {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	return template.Must(pgr.NewTemplate().Funcs(template.FuncMap{
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
	}).Parse(format))
}

func randIncBy(p *pgr.Bar, d time.Duration, seed int64) {
	r := rand.New(rand.NewSource(seed))
	for {
		p.Add(int64(r.Int() % 5))
		time.Sleep(d)
	}
}
