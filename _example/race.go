// +build ignore

package main

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
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

	poller := pgr.NewPoller(300 * time.Millisecond).
		Add(b1).
		Add(b2).
		Add(b3).
		Add(b4).
		Add(b5)

	now := time.Now().Unix()
	go randIncBy(b1, 300*time.Millisecond, now+1)
	go randIncBy(b2, 300*time.Millisecond, now+2)
	go randIncBy(b3, 300*time.Millisecond, now+3)
	go randIncBy(b4, 300*time.Millisecond, now+4)
	go randIncBy(b5, 300*time.Millisecond, now+5)

	poller.Show(context.Background())
}

var mu sync.Mutex
var no = make(map[string]int, 5)
var currentNo = 1

func makeBar(name, color string, width int) *pgr.Bar {
	prefix := fmt.Sprintf("({{ %s %q }})", color, name)
	pad := strings.Repeat(" ", width-runewidth.StringWidth(name))
	funcMap := funcMap()
	return pgr.NewBar(100, parseTemplate(funcMap, prefix+pad+` {{ rbar . "[" " " "🏇" "_" "]" 70 }}`)).
		OnFinishFunc(func(*pgr.Bar) string {
			mu.Lock()
			defer mu.Unlock()
			if _, exists := no[name]; !exists {
				no[name] = currentNo
				currentNo++
			}
			if f, ok := funcMap[color].(func(string) string); ok {
				return fmt.Sprintf(`(%s)%s Goal! (No.%d)`, f(name), pad, no[name])
			}
			return ""
		})
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
