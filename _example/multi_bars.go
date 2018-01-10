// +build ignore

package main

import (
	"context"
	"strings"
	"text/template"
	"time"

	"github.com/tyru/pgr"
)

func main() {
	p1 := pgr.NewBar("p1").SetTemplate(hasiruTemplate())
	p2 := pgr.NewBar("p2").SetTotal(200)
	p3 := pgr.NewBar("p3")
	go incBy(p1, 30*time.Millisecond)
	go incBy(p2, 20*time.Millisecond)
	go incBy(p3, 40*time.Millisecond)

	poller := pgr.NewPoller(100 * time.Millisecond).Add(p1)

	go func() {
		time.Sleep(1 * time.Second)
		poller.Add(p2)
		time.Sleep(1 * time.Second)
		poller.Add(p3)
	}()

	ctx := context.Background()
	// ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	// defer cancel()

	poller.Show(ctx)
}

func hasiruTemplate() *template.Template {
	t := pgr.NewBarTemplate().Funcs(template.FuncMap{
		"hasiru": func() func() string {
			forward := true
			dash := 0
			const maxDash = 10
			return func() string {
				var aa string
				if forward {
					if dash == 0 {
						aa = "┏( ^o^)┛"
					} else {
						aa = strings.Repeat("　", dash-1) + "三┏( ^o^)┛"
					}
				} else {
					if dash == 0 {
						aa = strings.Repeat("　", maxDash) + "┗(^o^ )┓"
					} else {
						aa = strings.Repeat("　", maxDash-dash) + "┗(^o^ )┓三"
					}
				}
				if dash >= maxDash {
					forward = !forward
					dash = 0
				} else {
					dash++
				}
				return aa
			}
		}(),
	})
	return template.Must(t.Parse(`{{ hasiru }}`))
}

func incBy(p *pgr.Bar, d time.Duration) {
	for {
		p.Inc()
		time.Sleep(d)
	}
}
